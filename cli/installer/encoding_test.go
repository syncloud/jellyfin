package installer

import (
	"os"
	"path"
	"strings"
	"testing"

	"go.uber.org/zap"
)

const encodingWithNilPreset = `<?xml version="1.0" encoding="utf-8"?>
<EncodingOptions xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <H264Crf>19</H264Crf>
  <EncoderPreset xsi:nil="true" />
  <DeinterlaceMethod>yadif</DeinterlaceMethod>
</EncodingOptions>
`

func TestDropNilEncoderPresetKeepsTheOtherSettings(t *testing.T) {
	file := writeEncoding(t, encodingWithNilPreset)

	err := DropNilEncoderPreset(file, zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	result := readEncoding(t, file)
	if strings.Contains(result, "EncoderPreset") {
		t.Fatalf("expected the nil EncoderPreset to be gone, got %q", result)
	}
	if !strings.Contains(result, "<H264Crf>19</H264Crf>") {
		t.Fatalf("expected H264Crf to survive, got %q", result)
	}
	if !strings.Contains(result, "<DeinterlaceMethod>yadif</DeinterlaceMethod>") {
		t.Fatalf("expected DeinterlaceMethod to survive, got %q", result)
	}
}

func TestDropNilEncoderPresetKeepsAnExplicitPreset(t *testing.T) {
	file := writeEncoding(t, strings.Replace(encodingWithNilPreset,
		`<EncoderPreset xsi:nil="true" />`, `<EncoderPreset>veryfast</EncoderPreset>`, 1))

	err := DropNilEncoderPreset(file, zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	result := readEncoding(t, file)
	if !strings.Contains(result, "<EncoderPreset>veryfast</EncoderPreset>") {
		t.Fatalf("expected the explicit preset to survive, got %q", result)
	}
}

func TestDropNilEncoderPresetKeepsTheFileMode(t *testing.T) {
	file := writeEncoding(t, encodingWithNilPreset)
	err := os.Chmod(file, 0640)
	if err != nil {
		t.Fatal(err)
	}

	err = DropNilEncoderPreset(file, zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0640 {
		t.Fatalf("expected the mode to be kept, got %v", info.Mode().Perm())
	}
}

func TestDropNilEncoderPresetWhenThereIsNoFile(t *testing.T) {
	err := DropNilEncoderPreset(path.Join(t.TempDir(), "encoding.xml"), zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}
}

func writeEncoding(t *testing.T, content string) string {
	t.Helper()
	file := path.Join(t.TempDir(), "encoding.xml")
	err := os.WriteFile(file, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func readEncoding(t *testing.T, file string) string {
	t.Helper()
	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
