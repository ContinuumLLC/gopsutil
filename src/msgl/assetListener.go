package msgl

import (
	"errors"

	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

//errors
const (
	ErrInvalidPluginPath = "ErrInvalidPluginPath"
)

//AssetListenerFactory contains function that returns RouteHandler
type AssetListenerFactory struct {
}

//GetAssetListener returns route handler
func (AssetListenerFactory) GetAssetListener(deps model.HandlerDependencies) model.AssetListener {
	config, _ := deps.GetConfigService(deps).GetAssetPluginConfig()
	logCfg := &logging.Config{
		AllowedLogLevel: logging.INFO,
		LogFileName:     deps.GetServiceInit().GetLogFilePath(),
		MaxFileSizeInMB: config.MaxLogFileSizeInMB,
		OldFileToKeep:   config.OldLogFileToKeep,
	}

	logCfg.SetLogLevel(config.LogLevel)
	logging.GetLoggerFactory().Update(*logCfg)

	return processAssetImpl{
		dependencies: deps,
		cfg:          config,
		logger:       logging.GetLoggerFactory().Get(),
	}
}

// HandleRoute is plugin path route handler
type HandleRoute func(req *protocol.Request) (res *protocol.Response, err error)

type processAssetImpl struct {
	server       protocol.Server
	dependencies model.HandlerDependencies
	cfg          *model.AssetPluginConfig
	logger       logging.Logger
}

func (pp processAssetImpl) registerRoutes() {
	pp.server.RegisterRoutes(
		&protocol.Route{Path: pp.cfg.PluginPath.AssetCollection,
			Handle: pp.dependencies.GetHandler(pp.dependencies, pp.cfg).HandleAsset},
		&protocol.Route{Path: pp.cfg.PluginPath.Configuration,
			Handle: pp.dependencies.GetHandler(pp.dependencies, pp.cfg).HandleConfig},
	)
}

func (pp processAssetImpl) Process() error {

	stdin := pp.dependencies.GetReader()
	stdout := pp.dependencies.GetWriter()

	pp.server = pp.dependencies.GetServer(stdin, stdout)
	pp.registerRoutes()
	var resp *protocol.Response
	request, err := pp.server.ReceiveRequest()
	if err != nil {
		err = exception.New(model.ErrAssetMsgListener, err)
		pp.sendErrorResponse(protocol.StatusCodeInternalError, resp, err)
		return err
	}

	route := request.Route
	if route == nil {
		err = errors.New(ErrInvalidPluginPath)
		pp.sendErrorResponse(protocol.PathNotFound, resp, err)
		return err
	}

	response, err := route.Handle(request)
	if err != nil {
		pp.sendErrorResponse(protocol.StatusCodeInternalError, resp, err)
		return err
	}

	pp.server.SendResponse(response)
	return nil
}

func (pp processAssetImpl) sendErrorResponse(code protocol.ResponseStatus, resp *protocol.Response, err error) {
	if resp == nil {
		resp = protocol.NewResponse()
	}
	pp.logger.Logf(logging.ERROR, "Response failure: Status %d", code)
	resp.SetError(code, err)
	pp.server.SendResponse(resp)
}
