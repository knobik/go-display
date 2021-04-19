// +build go1.10

package display

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation
#include <CoreGraphics/CoreGraphics.h>

void* CompatCGDisplayCreateImageForRect(CGDirectDisplayID display, CGRect rect) {
	return CGDisplayCreateImageForRect(display, rect);
}

void CompatCGImageRelease(void* image) {
	CGImageRelease(image);
}

void* CompatCGImageCreateCopyWithColorSpace(void* image, CGColorSpaceRef space) {
	return CGImageCreateCopyWithColorSpace((CGImageRef)image, space);
}

void CompatCGContextDrawImage(CGContextRef c, CGRect rect, void* image) {
	CGContextDrawImage(c, rect, (CGImageRef)image);
}
*/
import "C"

import (
	"unsafe"
)

func NumActiveDisplays() int {
	var count C.uint32_t = 0
	if C.CGGetActiveDisplayList(0, nil, &count) == C.kCGErrorSuccess {
		return int(count)
	} else {
		return 0
	}
}

func GetDisplayBounds(displayIndex int) Bounds {
	id := getDisplayId(displayIndex)
	main := C.CGMainDisplayID()

	var rect Bounds

	bounds := getCoreGraphicsCoordinateOfDisplay(id)
	rect.Left = int(bounds.origin.x)
	if main == id {
		rect.Top = 0
	} else {
		mainBounds := getCoreGraphicsCoordinateOfDisplay(main)
		mainHeight := mainBounds.size.height
		rect.Top = int(mainHeight - (bounds.origin.y + bounds.size.height))
	}
	rect.Right = rect.Left + int(bounds.size.width)
	rect.Bottom = rect.Top + int(bounds.size.height)

	return rect
}

func getDisplayId(displayIndex int) C.CGDirectDisplayID {
	main := C.CGMainDisplayID()
	if displayIndex == 0 {
		return main
	} else {
		n := NumActiveDisplays()
		ids := make([]C.CGDirectDisplayID, n)
		if C.CGGetActiveDisplayList(C.uint32_t(n), (*C.CGDirectDisplayID)(unsafe.Pointer(&ids[0])), nil) != C.kCGErrorSuccess {
			return 0
		}
		index := 0
		for i := 0; i < n; i++ {
			if ids[i] == main {
				continue
			}
			index++
			if index == displayIndex {
				return ids[i]
			}
		}
	}

	return 0
}

func getCoreGraphicsCoordinateOfDisplay(id C.CGDirectDisplayID) C.CGRect {
	main := C.CGDisplayBounds(C.CGMainDisplayID())
	r := C.CGDisplayBounds(id)
	return C.CGRectMake(r.origin.x, -r.origin.y-r.size.height+main.size.height,
		r.size.width, r.size.height)
}

func getCoreGraphicsCoordinateFromWindowsCoordinate(p C.CGPoint, mainDisplayBounds C.CGRect) C.CGPoint {
	return C.CGPointMake(p.x, mainDisplayBounds.size.height-p.y)
}

func createBitmapContext(width int, height int, data *C.uint32_t, bytesPerRow int) C.CGContextRef {
	colorSpace := createColorspace()
	if colorSpace == 0 {
		return 0
	}
	defer C.CGColorSpaceRelease(colorSpace)

	return C.CGBitmapContextCreate(unsafe.Pointer(data),
		C.size_t(width),
		C.size_t(height),
		8, // bits per component
		C.size_t(bytesPerRow),
		colorSpace,
		C.kCGImageAlphaNoneSkipFirst)
}

func createColorspace() C.CGColorSpaceRef {
	return C.CGColorSpaceCreateWithName(C.kCGColorSpaceSRGB)
}

func activeDisplayList() []C.CGDirectDisplayID {
	count := C.uint32_t(NumActiveDisplays())
	ret := make([]C.CGDirectDisplayID, count)
	if count > 0 && C.CGGetActiveDisplayList(count, (*C.CGDirectDisplayID)(unsafe.Pointer(&ret[0])), nil) == C.kCGErrorSuccess {
		return ret
	} else {
		return make([]C.CGDirectDisplayID, 0)
	}
}
