package doctron_core

import (
	"io/ioutil"
	"os"
	"os/exec"
)

type html2svg struct {
	*html2pdf
}

func (ins *html2svg) Convert() ([]byte, error) {
	buf, err := ins.html2pdf.Convert()
	ins.Log("pdf len:[%d]", len(buf))

	pdfFileName, err := ins.saveToPdf(buf)
	if err != nil {
		return nil, err
	}

	defer os.Remove(pdfFileName)

	svgFile, err := os.CreateTemp("", "pdf2svg_*.svg")
	if err != nil {
		return nil, err
	}
	svgFileName := svgFile.Name()
	//defer os.Remove(svgFileName)
	svgFile.Close()
	ins.Log("svg file:[%s]", svgFileName)
	output, err := exec.Command("pdf2svg", pdfFileName, svgFileName).Output()
	if err != nil {
		return nil, err
	}

	ins.Log("command output:[%s]", output)
	ins.buf, err = ioutil.ReadFile(svgFileName)
	if err != nil {
		return nil, err
	}

	//reader := bytes.NewReader(buf)
	//ins.buf, err = pdf2svg.Convert(reader)

	//if err != nil {
	//	return nil, err
	//}

	return ins.buf, err
}

func (ins *html2svg) saveToPdf(buf []byte) (string, error) {
	pdfFile, err := os.CreateTemp("", "pdf2svg_*.pdf")
	if err != nil {
		return "", err
	}

	defer pdfFile.Close()

	ins.Log("pdf filename:[%s]", pdfFile.Name())

	_, err = pdfFile.Write(buf)
	if err != nil {
		return pdfFile.Name(), err
	}

	return pdfFile.Name(), nil
}
