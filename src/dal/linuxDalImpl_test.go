package dal

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"strings"

	"fmt"

	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/asset"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model"
	"github.com/ContinuumLLC/platform-asset-plugin/src/model/mock"
	eMock "github.com/ContinuumLLC/platform-common-lib/src/env/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	pMock "github.com/ContinuumLLC/platform-common-lib/src/procParser/mock"
	"github.com/golang/mock/gomock"
)

var diskData = &procParser.Data{
	Lines: []procParser.Line{
		procParser.Line{
			Values: []string{"11", "0", "2346782340", "sda"},
		},
		procParser.Line{
			Values: []string{"11", "0", "35345", "sda1"},
		},
		procParser.Line{
			Values: []string{"11", "0", "5646", "sda2"},
		},
		procParser.Line{
			Values: []string{"11", "0", "6757", "sda5"},
		},
		procParser.Line{
			Values: []string{"11", "0", "345345", "sr0"},
		},
		procParser.Line{
			Values: []string{"11", "0", "345345", "xyz"},
		},
	},
}

const (
	hwXML  string = "<list></list>"
	hwXML1 string = `<?xml version="1.0" standalone="yes" ?>
<!-- generated by lshw-B.02.17 -->
<!-- GCC 5.3.1 20160413 -->
<!-- Linux 4.4.0-66-generic #87-Ubuntu SMP Fri Mar 3 15:29:05 UTC 2017 x86_64 -->
<!-- GNU libc 2 (glibc 2.23) -->
<list>
<node id="milinda-virtualbox" claimed="true" class="system" handle="DMI:0001">
 <description>Computer</description>
 <product>VirtualBox</product>
 <vendor>innotek GmbH</vendor>
 <version>1.2</version>
 <serial>0</serial>
 <width units="bits">64</width>
 <configuration>
  <setting id="family" value="Virtual Machine" />
  <setting id="uuid" value="E8F8DC61-D6D2-413D-86FE-A0E3E9807FC2" />
 </configuration>
 <capabilities>
  <capability id="smbios-2.5" >SMBIOS version 2.5</capability>
  <capability id="dmi-2.5" >DMI version 2.5</capability>
  <capability id="vsyscall32" >32-bit processes</capability>
 </capabilities>
  <node id="core" claimed="true" class="bus" handle="DMI:0008">
   <description>Motherboard</description>
   <product>VirtualBox</product>
   <vendor>Oracle Corporation</vendor>
   <physid>0</physid>
   <version>1.2</version>
   <serial>0</serial>
    <node id="firmware" claimed="true" class="memory" handle="">
     <description>BIOS</description>
     <vendor>innotek GmbH</vendor>
     <physid>0</physid>
     <version>VirtualBox</version>
     <date>12/01/2006</date>
     <size units="bytes">131072</size>
     <capabilities>
      <capability id="isa" >ISA bus</capability>
      <capability id="pci" >PCI bus</capability>
      <capability id="cdboot" >Booting from CD-ROM/DVD</capability>
      <capability id="bootselect" >Selectable boot path</capability>
      <capability id="int9keyboard" >i8042 keyboard controller</capability>
      <capability id="int10video" >INT10 CGA/Mono video</capability>
      <capability id="acpi" >ACPI</capability>
     </capabilities>
    </node>
    <node id="memory" claimed="true" class="memory" handle="">
     <description>System memory</description>
     <physid>1</physid>
     <size units="bytes">10485252096</size>
    </node>
  <node id="network" claimed="true" class="network" handle="PCI:0000:00:03.0">
   <description>Ethernet interface</description>
   <product>82540EM Gigabit Ethernet Controller</product>
   <vendor>Intel Corporation</vendor>
   <physid>3</physid>
   <businfo>pci@0000:00:03.0</businfo>
   <logicalname>lo</logicalname>
   <version>02</version>
   <serial>08:00:27:57:ec:1b</serial>
   <size units="bit/s">1000000000</size>
   <capacity>1000000000</capacity>
   <width units="bits">32</width>
   <clock units="Hz">66000000</clock>
   <configuration>
    <setting id="autonegotiation" value="on" />
    <setting id="broadcast" value="yes" />
    <setting id="driver" value="e1000" />
    <setting id="driverversion" value="7.3.21-k8-NAPI" />
    <setting id="duplex" value="full" />
    <setting id="ip" value="10.0.2.15" />
    <setting id="latency" value="64" />
    <setting id="link" value="yes" />
    <setting id="mingnt" value="255" />
    <setting id="multicast" value="yes" />
    <setting id="port" value="twisted pair" />
    <setting id="speed" value="1Gbit/s" />
   </configuration>
   <capabilities>
    <capability id="pm" >Power Management</capability>
    <capability id="pcix" >PCI-X</capability>
    <capability id="bus_master" >bus mastering</capability>
    <capability id="cap_list" >PCI capabilities listing</capability>
    <capability id="ethernet" />
    <capability id="physical" >Physical interface</capability>
    <capability id="tp" >twisted pair</capability>
    <capability id="10bt" >10Mbit/s</capability>
    <capability id="10bt-fd" >10Mbit/s (full duplex)</capability>
    <capability id="100bt" >100Mbit/s</capability>
    <capability id="100bt-fd" >100Mbit/s (full duplex)</capability>
    <capability id="1000bt-fd" >1Gbit/s (full duplex)</capability>
    <capability id="autonegotiation" >Auto-negotiation</capability>
   </capabilities>
   <resources>
    <resource type="irq" value="19" />
    <resource type="memory" value="f0000000-f001ffff" />
    <resource type="ioport" value="d010(size=8)" />
   </resources>
  </node>
 <node id="network" claimed="true" class="network" handle="PCI:0000:00:03.0">
   <description>Ethernet interface</description>
   <product>82540EM Gigabit Ethernet Controller</product>
   <vendor>Intel Corporation</vendor>
   <physid>3</physid>
   <businfo>pci@0000:00:03.0</businfo>
   <logicalname>enp0s3</logicalname>
   <version>02</version>
   <serial>08:00:27:57:ec:1b</serial>
   <size units="bit/s">1000000000</size>
   <capacity>1000000000</capacity>
   <width units="bits">32</width>
   <clock units="Hz">66000000</clock>  
  </node> 
  <node id="usb" claimed="true" class="bus" handle="PCI:0000:00:06.0">
   <description>USB controller</description>
   <product>KeyLargo/Intrepid USB</product>
   <vendor>Apple Inc.</vendor>
   <physid>6</physid>
   <businfo>pci@0000:00:06.0</businfo>
   <version>00</version>
   <width units="bits">32</width>
   <clock units="Hz">33000000</clock>
   <configuration>
    <setting id="driver" value="ohci-pci" />
    <setting id="latency" value="64" />
   </configuration>
   <capabilities>
    <capability id="ohci" >Open Host Controller Interface</capability>
    <capability id="bus_master" >bus mastering</capability>
    <capability id="cap_list" >PCI capabilities listing</capability>
   </capabilities>
   <resources>
    <resource type="irq" value="22" />
    <resource type="memory" value="f0804000-f0804fff" />
   </resources>
    <node id="usbhost" claimed="true" class="bus" handle="USB:1:1">
     <product>OHCI PCI host controller</product>
     <vendor>Linux 4.4.0-66-generic ohci_hcd</vendor>
     <physid>1</physid>
     <businfo>usb@1</businfo>
     <logicalname>usb1</logicalname>
     <version>4.04</version>
     <configuration>
      <setting id="driver" value="hub" />
      <setting id="slots" value="12" />
      <setting id="speed" value="12Mbit/s" />
     </configuration>
     <capabilities>
      <capability id="usb-1.10" >USB 1.1</capability>
     </capabilities>
    </node>
  </node>
  <node id="cdrom" claimed="true" class="disk" handle="SCSI:01:00:00:00">
   <description>DVD reader</description>
   <physid>0.0.0</physid>
   <businfo>scsi@1:0.0.0</businfo>
   <logicalname>/dev/cdrom</logicalname>
   <logicalname>/dev/dvd</logicalname>
   <logicalname>/dev/sr0</logicalname>
   <dev>11:0</dev>
   <configuration>
    <setting id="status" value="nodisc" />
   </configuration>
   <capabilities>
    <capability id="audio" >Audio CD playback</capability>
    <capability id="dvd" >DVD playback</capability>
   </capabilities>
  </node>
  <node id="disk" claimed="true" class="disk" handle="SCSI:02:00:00:00">
   <description>ATA Disk</description>
   <product>VBOX HARDDISK</product>
   <physid>0.0.0</physid>
   <businfo>scsi@2:0.0.0</businfo>
   <logicalname>/dev/sda</logicalname>
   <dev>8:0</dev>
   <version>1.0</version>
   <serial>VB7f4a1ba4-6ef7655d</serial>
   <size units="bytes">64424509440</size>
   <configuration>
    <setting id="ansiversion" value="5" />
    <setting id="logicalsectorsize" value="512" />
    <setting id="sectorsize" value="512" />
    <setting id="signature" value="29043892" />
   </configuration>
   <capabilities>
    <capability id="partitioned" >Partitioned disk</capability>
    <capability id="partitioned:dos" >MS-DOS partition table</capability>
   </capabilities>
    <node id="volume:0" claimed="true" class="volume" handle="">
     <description>EXT4 volume</description>
     <vendor>Linux</vendor>
     <physid>1</physid>
     <businfo>scsi@2:0.0.0,1</businfo>
     <logicalname>/dev/sda1</logicalname>
     <logicalname>/</logicalname>
     <dev>8:1</dev>
     <version>1.0</version>
     <serial>ecc36b30-d5d2-40a0-9962-88661930be29</serial>
     <size units="bytes">55833526272</size>
     <capacity>55833526272</capacity>
     <configuration>
      <setting id="created" value="2016-12-15 13:38:16" />
      <setting id="filesystem" value="ext4" />
      <setting id="lastmountpoint" value="/" />
      <setting id="modified" value="2017-03-19 20:50:51" />
      <setting id="mount.fstype" value="ext4" />
      <setting id="mount.options" value="rw,relatime,errors=remount-ro,data=ordered" />
      <setting id="mounted" value="2017-03-18 11:28:22" />
      <setting id="state" value="mounted" />
     </configuration>
     <capabilities>
      <capability id="primary" >Primary partition</capability>
      <capability id="bootable" >Bootable partition (active)</capability>
      <capability id="journaled" />
      <capability id="extended_attributes" >Extended Attributes</capability>
      <capability id="large_files" >4GB+ files</capability>
      <capability id="huge_files" >16TB+ files</capability>
      <capability id="dir_nlink" >directories with 65000+ subdirs</capability>
      <capability id="extents" >extent-based allocation</capability>
      <capability id="ext4" />
      <capability id="ext2" >EXT2/EXT3</capability>
      <capability id="initialized" >initialized volume</capability>
     </capabilities>
    </node>
    <node id="volume:1" claimed="true" class="volume" handle="">
     <description>Extended partition</description>
     <physid>2</physid>
     <businfo>scsi@2:0.0.0,2</businfo>
     <logicalname>/dev/sda2</logicalname>
     <dev>8:2</dev>
     <size units="bytes">8587838464</size>
     <capacity>8587838464</capacity>
     <capabilities>
      <capability id="primary" >Primary partition</capability>
      <capability id="extended" >Extended partition</capability>
      <capability id="partitioned" >Partitioned disk</capability>
      <capability id="partitioned:extended" >Extended partition</capability>
     </capabilities>
      <node id="logicalvolume" claimed="true" class="volume" handle="">
       <description>Linux swap / Solaris partition</description>
       <physid>5</physid>
       <logicalname>/dev/sda5</logicalname>
       <dev>8:5</dev>
       <capacity>8587837440</capacity>
       <capabilities>
        <capability id="nofs" >No filesystem</capability>
       </capabilities>
      </node>
    </node>
  </node>
  </node>
</node>
</list>`
)

func setupGetCommandReader(t *testing.T, parseErr error, commandReaderErr error) (*gomock.Controller, *mock.MockAssetDalDependencies) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any()).Return(reader, commandReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if commandReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(&procParser.Data{}, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

//TODO - Duplicate function as setupGetCommandReader. Need to relook at it.
func setupGetCommandReader2(t *testing.T, parseErr error, commandReaderErr error, data *procParser.Data) (*gomock.Controller, *mock.MockAssetDalDependencies) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, commandReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if commandReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(data, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

func setupGetFileReader(t *testing.T, parseErr error, fileReaderErr error, parseData *procParser.Data) (*gomock.Controller, *mock.MockAssetDalDependencies) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, fileReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if fileReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(parseData, parseErr)
	}
	mockAssetDalD.EXPECT().GetParser().Return(mockParser)

	return ctrl, mockAssetDalD
}

func setupAddGetFileReader(ctrl *gomock.Controller, mockAssetDalD *mock.MockAssetDalDependencies, parseErr error, fileReaderErr error) {
	mockEnv := eMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, fileReaderErr)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv)

	mockParser := pMock.NewMockParser(ctrl)
	if fileReaderErr == nil {
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(&procParser.Data{}, parseErr)
	}
}

func TestGetOSCommandErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, errors.New(model.ErrExecuteCommandFailed))
	defer ctrl.Finish()

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetOSInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetOSFileErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, nil)
	defer ctrl.Finish()

	setupAddGetFileReader(ctrl, mockAssetDalD, nil, errors.New(model.ErrFileReadFailed))

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetOSInfo()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrFileReadFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrFileReadFailed, err)
	}
}

// TODO - fix error
// func TestGetOSNoErr(t *testing.T) {
// 	ctrl, mockAssetDalD := setupGetCommandReader(t, nil, nil)
// 	defer ctrl.Finish()

// 	setupAddGetFileReader(ctrl, mockAssetDalD, nil, nil)

//
//
// 	_, err := assetDalImpl{
// 		Factory: mockAssetDalD,
// 		Logger : logging.GetLoggerFactory().Get(),
// 	}.GetOS()
// 	if err != nil {
// 		t.Errorf("Unexpected error : %v", err)
// 	}
// }

func setupGetSystemInfo(t *testing.T, times int, err error) (*gomock.Controller, error) {
	ctrl := gomock.NewController(t)

	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	var str string
	switch times {
	case 1:
		str = cSysProductCmd
	case 2:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		str = cSysTz
	case 3:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTz).Return("", nil)
		str = cSysTzd
	case 4:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTz).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTzd).Return("", nil)
		str = cSysSerialNo
	case 5:
		mockEnv.EXPECT().ExecuteBash(cSysProductCmd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTz).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysTzd).Return("", nil)
		mockEnv.EXPECT().ExecuteBash(cSysSerialNo).Return("", nil)
		str = cSysHostname

	}
	mockEnv.EXPECT().ExecuteBash(str).Return("", err)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(times)

	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetSystemInfo()
	return ctrl, e
}

func TestGetSystemInfoErr(t *testing.T) {
	cmdExeArr := []int{1, 2, 3, 4, 5}
	for _, i := range cmdExeArr {
		ctrl, err := setupGetSystemInfo(t, i, errors.New(model.ErrExecuteCommandFailed))
		defer ctrl.Finish()
		if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
			t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
		}
	}
}

func TestGetSystemNoErr(t *testing.T) {
	ctrl, err := setupGetSystemInfo(t, 5, nil)
	defer ctrl.Finish()
	if err != nil {
		t.Errorf("Unexpected error received  : %v", err)
	}
}

func TestGetMemoryInfoErr(t *testing.T) {
	parseError := model.ErrFileReadFailed
	_, mockAssetDalD := setupGetFileReader(t, errors.New(parseError), nil, nil)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetMemoryInfo()

	if err == nil || err.Error() != parseError {
		t.Error("Error expected but not returned")
	}
}

func TestGetDataFromMap(t *testing.T) {
	data := procParser.Data{
		Map: make(map[string]procParser.Line, 1),
	}
	data.Map["MemTotal"] = procParser.Line{Values: []string{"MemTotal", "InvalidNumber", "KB"}}

	util := dalUtil{}
	val := util.getDataFromMap("MemTotal", &data)

	if val != 0 {
		t.Errorf("Expected 0, returned %d", val)
	}
}

func TestGetMemoryInfoNoErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := procParser.Data{
		Map: make(map[string]procParser.Line, 1),
	}
	data.Map["MemTotal"] = procParser.Line{Values: []string{"physicalTotalBytes", "1", "KB"}}

	_, mockAssetDalD := setupGetFileReader(t, nil, nil, &data)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetMemoryInfo()

	if err != nil {
		t.Errorf("Unexpected error received  : %v", err)
	}
}

func TestGetProcessorInfoErr(t *testing.T) {
	parseError := model.ErrFileReadFailed
	ctrl, mockAssetDalD := setupGetFileReader(t, errors.New(parseError), nil, nil)
	defer ctrl.Finish()

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetProcessorInfo()

	if err == nil || err.Error() != parseError {
		t.Error("Error expected but not returned")
	}
}

func TestGetProcessorInfoBashErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetFileReader(t, nil, nil, nil)
	defer ctrl.Finish()
	envMock := eMock.NewMockEnv(ctrl)
	envMock.EXPECT().ExecuteBash(cCPUArcCmd).Return("", errors.New(model.ErrExecuteCommandFailed))
	mockAssetDalD.EXPECT().GetEnv().Return(envMock)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetProcessorInfo()

	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetProcessorInfoNoErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetFileReader(t, nil, nil, &procParser.Data{})
	defer ctrl.Finish()
	envMock := eMock.NewMockEnv(ctrl)
	envMock.EXPECT().ExecuteBash(cCPUArcCmd).Return("", nil)
	mockAssetDalD.EXPECT().GetEnv().Return(envMock)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetProcessorInfo()

	if err != nil {
		t.Errorf("Unexpected error returned : %v", err)
	}
}

func setupEnv(t *testing.T) (*gomock.Controller, *mock.MockAssetDalDependencies, *eMock.MockEnv) {
	ctrl := gomock.NewController(t)
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(1)
	v = nil

	return ctrl, mockAssetDalD, mockEnv
}

func TestReadHwList(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return("<list></list>", nil)

	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.readHwList()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestReadHwListError(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return("<list></list>", errors.New("readHwListErr"))

	v = nil
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.readHwList()
	if e == nil || e.Error() != model.ErrExecuteCommandFailed {
		t.Errorf("Expecting model.ErrExecuteCommandFailed , Unexpected error %v", e)
	}
}

func TestReadHwListErr(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()
	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return("nu$756ll", nil)

	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.readHwList()
	if e == nil {
		t.Error("Expecting EOF error ")
	}
}

func TestGetNetworkInfoNew(t *testing.T) {
	fmt.Println("TestGetNetworkInfoNew")
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	// v := List{}
	// xml.Unmarshal([]byte(hwXML), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetNetworkInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetNetworkInfoNew2(t *testing.T) {
	fmt.Println("TestGetNetworkInfoNew2")
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, nil)
	// v := List{}
	// xml.Unmarshal([]byte(hwXML), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetNetworkInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBiosInfo(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	// v := List{}
	// xml.Unmarshal([]byte(hwXML), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetBiosInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBiosInfo2(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, nil)
	// v := List{}
	// xml.Unmarshal([]byte(hwXML1), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetBiosInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBiosInfoError(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, errors.New("XMLReadError"))
	// v := List{}
	// xml.Unmarshal([]byte(hwXML1), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetBiosInfo()
	if e == nil || e.Error() != model.ErrExecuteCommandFailed {
		t.Errorf("Unexpected error %v, was expecting model.ErrExecuteCommandFailed", e)
	}
}

func TestGetBaseBoardInfo(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetBaseBoardInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBaseBoardInfo2(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, nil)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetBaseBoardInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetBaseBoardInfoError(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, errors.New("XMLReadErr"))
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetBaseBoardInfo()
	if e == nil || e.Error() != model.ErrExecuteCommandFailed {
		t.Errorf("Unexpected error %v, was expecting model.ErrExecuteCommandFailed ", e)
	}
}

func TestGetDrivesInfo(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, nil)
	mockParser := pMock.NewMockParser(ctrl)

	mockAssetDalD.EXPECT().GetParser().Return(mockParser)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(1)
	mockParser.EXPECT().Parse(gomock.Any(), gomock.Any()).Return(diskData, nil)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetDrivesInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}
func TestGetDrivesInfo2(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, nil)
	// v := List{}
	// xml.Unmarshal([]byte(hwXML1), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetDrivesInfo()
	if e != nil {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetDrivesInfoError2(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, nil)
	mockParser := pMock.NewMockParser(ctrl)

	mockAssetDalD.EXPECT().GetParser().Return(mockParser)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(1)
	mockParser.EXPECT().Parse(gomock.Any(), gomock.Any()).Return(diskData, errors.New("ParseErr"))
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetDrivesInfo()
	if e == nil || e.Error() != "ParseErr" {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetDrivesInfoError3(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML, nil)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, errors.New("ReadErr"))
	//mockParser := pMock.NewMockParser(ctrl)

	//mockAssetDalD.EXPECT().GetParser().Return(mockParser)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(1)
	//mockParser.EXPECT().Parse(gomock.Any(), gomock.Any()).Return(diskData, errors.New("ParseErr"))
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetDrivesInfo()
	if e == nil || e.Error() != "ReadErr" {
		t.Errorf("Unexpected error %v", e)
	}
}

func TestGetDrivesInfoError(t *testing.T) {
	ctrl, mockAssetDalD, mockEnv := setupEnv(t)
	defer ctrl.Finish()

	mockEnv.EXPECT().ExecuteBash(cListHwAsXML).Return(hwXML1, errors.New("XMLReadErr"))
	// v := List{}
	// xml.Unmarshal([]byte(hwXML1), &v)
	_, e := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetDrivesInfo()
	if e == nil || e.Error() != model.ErrExecuteCommandFailed {
		t.Errorf("Unexpected error %v, was expecting model.ErrExecuteCommandFailed", e)
	}
}

func TestMapToArr(t *testing.T) {
	m := map[string]asset.AssetNetwork{
		"eth0": asset.AssetNetwork{},
		"eth1": asset.AssetNetwork{},
	}
	nArr := mapToArr(m)
	if l := len(nArr); l != 2 {
		t.Errorf("Expected length is %d but received %d", 2, l)
	}
}

func TestSetValnmcli(t *testing.T) {
	networks := map[string]asset.AssetNetwork{
		"eth0": asset.AssetNetwork{},
		"eth1": asset.AssetNetwork{},
	}
	mapArr := map[string]map[string][]string{
		"eth0": {
			"DHCP4.OPTION[11]":         []string{"DHCP4.OPTION[11]", "dhcp_server_identifier = 10.0.3.2"},
			"DHCP4.OPTION[9]":          []string{"DHCP4.OPTION[9]", "domain_name_servers = 10.2.17.6 10.2.17.25 10.2.17.17"},
			"DHCP4.OPTION[6]":          []string{"DHCP4.OPTION[6]", "ip_address = 10.0.3.15"},
			"DHCP4.OPTION[5]":          []string{"DHCP4.OPTION[5]", "ip_address 10.0.3.15"},
			"DHCP4.OPTION[7]":          []string{"DHCP4.OPTION[7]", "subnet_mask = 255.255.255.0"},
			"GENERAL.HWADDR":           []string{"GENERAL.HWADDR", "08:00:27:09:C7:82"},
			"GENERAL.FIRMWARE-VERSION": []string{"GENERAL.FIRMWARE-VERSION"},
		},
	}
	setValnmcli(networks, mapArr)
	if d := networks["eth0"].DhcpServer; d != "10.0.3.2" {
		t.Errorf("Expected value is 10.0.3.2 but received %s", d)
	}
	if i := networks["eth0"].IPv4; i != "10.0.3.15" {
		t.Errorf("Expected value is 10.0.3.15 but received %s", i)
	}
}

func TestGetNetworkInfo1(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader2(t, nil, errors.New(model.ErrExecuteCommandFailed), &procParser.Data{})
	defer ctrl.Finish()

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetNetworkInfo1()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetNetworkInfo1CommandDataErr(t *testing.T) {
	ctrl, mockAssetDalD := setupGetCommandReader2(t, nil, nil, &procParser.Data{
		Lines: []procParser.Line{
			procParser.Line{
				Values: []string{"*-network", "0"},
			},
			procParser.Line{
				Values: []string{"*product", "82540EM Gigabit Ethernet Controller"},
			},
			procParser.Line{
				Values: []string{"*-network", "1"},
			},
		},
	})

	mockEnv := eMock.NewMockEnv(ctrl)
	mockAssetDalD.EXPECT().GetEnv().Return(mockEnv).Times(1)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New(model.ErrExecuteCommandFailed))

	defer ctrl.Finish()

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetNetworkInfo1()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}

func TestGetNetworkInfo1CommandDataErr1(t *testing.T) {
	data := &procParser.Data{
		Lines: []procParser.Line{
			procParser.Line{
				Values: []string{"*-network", "0"},
			},
			procParser.Line{
				Values: []string{"*product", "82540EM Gigabit Ethernet Controller"},
			},
			procParser.Line{
				Values: []string{"logical name", "enp0s3"},
			},
			procParser.Line{
				Values: []string{"*-network", "1"},
			},
			procParser.Line{
				Values: []string{"logical name", "enp0s4"},
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAssetDalD := mock.NewMockAssetDalDependencies(ctrl)

	mockEnv := eMock.NewMockEnv(ctrl)
	mockParser := pMock.NewMockParser(ctrl)
	byteReader := bytes.NewReader([]byte("data"))
	reader := ioutil.NopCloser(byteReader)
	gomock.InOrder(
		mockAssetDalD.EXPECT().GetParser().Return(mockParser),
		mockAssetDalD.EXPECT().GetEnv().Return(mockEnv),
		mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, nil),
		mockParser.EXPECT().Parse(gomock.Any(), reader).Return(data, nil),

		mockAssetDalD.EXPECT().GetEnv().Return(mockEnv),
		mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, errors.New("Err")),
	)

	_, err := assetDalImpl{
		Factory: mockAssetDalD,
		Logger:  logging.GetLoggerFactory().Get(),
	}.GetNetworkInfo1()
	if err == nil || !strings.HasPrefix(err.Error(), model.ErrExecuteCommandFailed) {
		t.Errorf("Expected error is %s, but received %v", model.ErrExecuteCommandFailed, err)
	}
}
