package streamdeck

import (
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	ELGATO_VID                         = "0fd9"
	USB_PID_STREAMDECK_MINI            = "0063"
	USB_PID_STREAMDECK_MINI_MK2        = "0090"
	USB_PID_STREAMDECK_MINI_MK2_MODULE = "00b8"
	USB_PID_STREAMDECK_MK2             = "0080"
	USB_PID_STREAMDECK_MK2_MODULE      = "00b9"
	USB_PID_STREAMDECK_MK2_SCISSOR     = "00a5"
	USB_PID_STREAMDECK_MK2_V2          = "00B9"
	USB_PID_STREAMDECK_NEO             = "009a"
	USB_PID_STREAMDECK_ORIGINAL        = "0060"
	USB_PID_STREAMDECK_ORIGINAL_V2     = "006d"
	USB_PID_STREAMDECK_PEDAL           = "0086"
	USB_PID_STREAMDECK_PLUS            = "0084"
	USB_PID_STREAMDECK_XL              = "006c"
	USB_PID_STREAMDECK_XL_V2           = "008f"
	USB_PID_STREAMDECK_STUDIO          = "00aa"
	USB_PID_STREAMDECK_XL_V2_MODULE    = "00ba"
)

// Firmware command IDs.
//
//nolint:revive
var (
	c_REV1_FIRMWARE   = []byte{0x04}
	c_REV1_RESET      = []byte{0x0b, 0x63}
	c_REV1_BRIGHTNESS = []byte{0x05, 0x55, 0xaa, 0xd1, 0x01}

	c_REV2_FIRMWARE   = []byte{0x05}
	c_REV2_RESET      = []byte{0x03, 0x02}
	c_REV2_BRIGHTNESS = []byte{0x03, 0x08}
)

type Key struct {
	Index   uint8
	Pressed bool
}

type DeckDevice struct {
	SysFs        string
	ID           string
	Serial       string
	Manufacturer string
	Product      string
	Path         string

	Columns uint8
	Rows    uint8
	Keys    uint8
	Pixels  uint
	DPI     uint
	Padding uint

	featureReportSize   int
	firmwareOffset      int
	keyStateOffset      int
	translateKeyIndex   func(index, columns uint8) uint8
	imagePageSize       int
	imagePageHeaderSize int
	flipImage           func(image.Image) image.Image
	toImageFormat       func(image.Image) ([]byte, error)
	imagePageHeader     func(pageIndex int, keyIndex uint8, payloadLength int, lastPage bool) []byte

	getFirmwareCommand   []byte
	resetCommand         []byte
	setBrightnessCommand []byte

	keyStateLength int
	device         *os.File
}

func FindDevices() ([]DeckDevice, error) {
	sysUSBPath := "/sys/bus/usb/devices"
	result := []DeckDevice{}

	// 1. Read the directory
	files, err := os.ReadDir(sysUSBPath)
	if err != nil {
		fmt.Println("Error reading sysfs:", err)
		return result, err
	}

	for _, file := range files {

		fullPath := filepath.Join(sysUSBPath, file.Name())

		// 2. Read VID and PID files inside the device folder
		vid, _ := readFileValue(filepath.Join(fullPath, "idVendor"))
		pid, _ := readFileValue(filepath.Join(fullPath, "idProduct"))

		// Check if we found valid IDs
		if vid != ELGATO_VID || pid == "" {
			continue
		}
		var dev DeckDevice
		path := ""
		serial := ""

		switch pid {
		case USB_PID_STREAMDECK_ORIGINAL:
			dev = DeckDevice{
				SysFs:                fullPath,
				ID:                   path,
				Serial:               serial,
				Columns:              5,
				Rows:                 3,
				Keys:                 15,
				Pixels:               72,
				DPI:                  124,
				Padding:              16,
				featureReportSize:    17,
				firmwareOffset:       5,
				keyStateOffset:       1,
				translateKeyIndex:    translateRightToLeft,
				imagePageSize:        7819,
				imagePageHeaderSize:  16,
				imagePageHeader:      rev1ImagePageHeader,
				flipImage:            flipHorizontally,
				toImageFormat:        toBMP,
				getFirmwareCommand:   c_REV1_FIRMWARE,
				resetCommand:         c_REV1_RESET,
				setBrightnessCommand: c_REV1_BRIGHTNESS,
			}
		case USB_PID_STREAMDECK_MINI, USB_PID_STREAMDECK_MINI_MK2:
			dev = DeckDevice{
				SysFs:                fullPath,
				ID:                   path,
				Serial:               serial,
				Columns:              3,
				Rows:                 2,
				Keys:                 6,
				Pixels:               80,
				DPI:                  138,
				Padding:              16,
				featureReportSize:    17,
				firmwareOffset:       5,
				keyStateOffset:       1,
				translateKeyIndex:    identity,
				imagePageSize:        1024,
				imagePageHeaderSize:  16,
				imagePageHeader:      miniImagePageHeader,
				flipImage:            rotateCounterclockwise,
				toImageFormat:        toBMP,
				getFirmwareCommand:   c_REV1_FIRMWARE,
				resetCommand:         c_REV1_RESET,
				setBrightnessCommand: c_REV1_BRIGHTNESS,
			}
		case USB_PID_STREAMDECK_ORIGINAL_V2, USB_PID_STREAMDECK_MK2:
			dev = DeckDevice{
				SysFs:                fullPath,
				ID:                   path,
				Serial:               serial,
				Columns:              5,
				Rows:                 3,
				Keys:                 15,
				Pixels:               72,
				DPI:                  124,
				Padding:              16,
				featureReportSize:    32,
				firmwareOffset:       6,
				keyStateOffset:       4,
				translateKeyIndex:    identity,
				imagePageSize:        1024,
				imagePageHeaderSize:  8,
				imagePageHeader:      rev2ImagePageHeader,
				flipImage:            flipHorizontallyAndVertically,
				toImageFormat:        toJPEG,
				getFirmwareCommand:   c_REV2_FIRMWARE,
				resetCommand:         c_REV2_RESET,
				setBrightnessCommand: c_REV2_BRIGHTNESS,
			}
		case USB_PID_STREAMDECK_XL:
			dev = DeckDevice{
				SysFs:                fullPath,
				ID:                   path,
				Serial:               serial,
				Columns:              8,
				Rows:                 4,
				Keys:                 32,
				Pixels:               96,
				DPI:                  166,
				Padding:              16,
				featureReportSize:    32,
				firmwareOffset:       6,
				keyStateOffset:       4,
				translateKeyIndex:    identity,
				imagePageSize:        1024,
				imagePageHeaderSize:  8,
				imagePageHeader:      rev2ImagePageHeader,
				flipImage:            flipHorizontallyAndVertically,
				toImageFormat:        toJPEG,
				getFirmwareCommand:   c_REV2_FIRMWARE,
				resetCommand:         c_REV2_RESET,
				setBrightnessCommand: c_REV2_BRIGHTNESS,
			}
		case USB_PID_STREAMDECK_NEO:
			dev = DeckDevice{
				SysFs:                fullPath,
				ID:                   path,
				Serial:               serial,
				Columns:              4,
				Rows:                 3,
				Keys:                 10,
				Pixels:               96,
				DPI:                  166,
				Padding:              16,
				featureReportSize:    32,
				firmwareOffset:       6,
				keyStateOffset:       4,
				translateKeyIndex:    identity,
				imagePageSize:        1024,
				imagePageHeaderSize:  8,
				imagePageHeader:      rev2ImagePageHeader,
				flipImage:            flipHorizontallyAndVertically,
				toImageFormat:        toJPEG,
				getFirmwareCommand:   c_REV2_FIRMWARE,
				resetCommand:         c_REV2_RESET,
				setBrightnessCommand: c_REV2_BRIGHTNESS,
			}
		}
		if dev.SysFs != "" {
			manufacturer, _ := readFileValue(filepath.Join(fullPath, "manufacturer"))
			product, _ := readFileValue(filepath.Join(fullPath, "product"))
			dev.Manufacturer = manufacturer
			dev.Product = product
			serial, _ = readFileValue(filepath.Join(fullPath, "serial"))
			dev.Serial = serial
			eventPath, err := findDevPath(fullPath)
			dev.Path = eventPath
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return result, err
			}

			fmt.Printf("USB Device: %s\n", fullPath)
			fmt.Printf("Input Path: %s\n", eventPath)
			fmt.Printf("   VID: %s\n", vid)
			fmt.Printf("   PID: %s\n", pid)
			dev.keyStateLength = int(dev.Columns * dev.Rows)
			result = append(result, dev)
		}
	}
	return result, nil
}

// Helper to read a one-line file and trim whitespace
func readFileValue(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func findDevPath(usbSysPath string) (string, error) {
	var eventNode string

	realPath, err := filepath.EvalSymlinks(usbSysPath)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate symlinks: %w", err)
	}

	err = filepath.WalkDir(realPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), "hidraw") {
			// Verify it is inside an "hidraw" subdirectory to avoid false positives
			if strings.Contains(path, "/hidraw/") {
				eventNode = d.Name()
				return filepath.SkipAll // Stop searching once found
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to walk sysfs: %w", err)
	}

	if eventNode == "" {
		return "", fmt.Errorf("no input event node found for %s (device might not be an input device)", usbSysPath)
	}

	return fmt.Sprintf("/dev/%s", eventNode), nil
}

func (dd *DeckDevice) OpenDeckDevice() error {
	// 1. Open the device file
	file, err := os.OpenFile(dd.Path, os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("❌ Error opening %s: %v (Requires root/input group permissions)\n", dd.Path, err)
		return err
	}
	dd.device = file
	return nil
}

func (dd *DeckDevice) Close() error {
	if dd.device != nil {
		err := dd.device.Close()
		dd.device = nil
		return err
	}
	return nil
}

func (dd *DeckDevice) ListenKeys() (chan Key, error) {
	if dd.device == nil {
		return nil, fmt.Errorf("device not opened")
	}

	kch := make(chan Key)
	keyBufferLen := dd.keyStateOffset + dd.keyStateLength
	oldKeyBuffer := make([]byte, keyBufferLen)
	go func() {
		for {

			keyBuffer := make([]byte, keyBufferLen)
			if _, err := dd.device.Read(keyBuffer); err != nil {
				close(kch)
				return
			}

			if len(keyBuffer) < keyBufferLen {
				continue
			}

			for i := dd.keyStateOffset; i < keyBufferLen; i++ {
				if keyBuffer[i] == oldKeyBuffer[i] {
					continue
				}
				keyIndex := uint8(i - dd.keyStateOffset)
				kch <- Key{
					Index:   dd.translateKeyIndex(keyIndex, dd.Columns),
					Pressed: keyBuffer[i] == 1,
				}
			}

			oldKeyBuffer = keyBuffer
		}
	}()

	return kch, nil
}

func (dd *DeckDevice) Write(data []byte) error {
	if dd.device == nil {
		return fmt.Errorf("device not opened")
	}
	_, err := dd.device.Write(data)
	return err
}

func (dd *DeckDevice) SetBrightness(percent uint8) error {
	if percent > 100 {
		percent = 100
	}

	report := make([]byte, len(dd.setBrightnessCommand)+1)
	copy(report, dd.setBrightnessCommand)
	report[len(report)-1] = percent

	return nil
}

// translateRightToLeft translates the given key index from right-to-left to
// left-to-right, based on the given number of columns.
func translateRightToLeft(index, columns uint8) uint8 {
	keyCol := index % columns
	return (index - keyCol) + (columns - 1) - keyCol
}

// identity returns the given key index as it is.
func identity(index, _ uint8) uint8 {
	return index
}

// // getFeatureReport from the device without worries about the correct payload size.
// func (dd *DeckDevice) getFeatureReport(id byte) ([]byte, error) {
// 	return dd.device.GetFeatureReport(id)
// }

// // setFeatureReport to the device without worries about the correct payload size.
// func (dd *DeckDevice) setFeatureReport(payload []byte) error {
// 	b := make([]byte, dd.featureReportSize-1)
// 	copy(b, payload[1:])
// 	return dd.device.SetFeatureReport(payload)
// }