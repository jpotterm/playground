package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

func main() {
	t := time.Date(2020, 10, 3, 0, 0, 0, 0, time.UTC)

	for _, image := range images {
		Process(image, t)
		t = t.Add(time.Minute)
	}
}

func Process(filename string, timestamp time.Time) {
	filepath := "images/in/" + filename

	jmp := jpegstructure.NewJpegMediaParser()

	intfc, err := jmp.ParseFile(filepath)
	if err != nil {
		panic(err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Update the CameraOwnerName tag.

	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		panic(err)
	}

	ifdPath := "IFD/Exif"

	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)
	if err != nil {
		panic(err)
	}

	updatedTimestampPhrase := exifcommon.ExifFullTimestampString(timestamp)

	fmt.Println(updatedTimestampPhrase)

	err = ifdIb.SetStandardWithName("DateTimeOriginal", updatedTimestampPhrase)
	if err != nil {
		panic(err)
	}

	// Update the exif segment.

	err = sl.SetExif(rootIb)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("images/out/" + filename)
	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	err = sl.Write(outputFile)
	if err != nil {
		panic(err)
	}
}
