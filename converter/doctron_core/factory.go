package doctron_core

import (
	"context"

	"github.com/lampnick/doctron/converter"
)

const (
	DoctronHtml2Pdf     = 1
	DoctronHtml2Svg     = 2
	DoctronHtml2Image   = 3
	DoctronPdf2Image    = 4
	DoctronPdfWatermark = 5
)

// new doctron
func NewDoctron(config converter.DoctronConfig) DoctronI {
	switch config.DoctronType {
	case DoctronHtml2Pdf:
		fac := &(html2PdfFactory{})
		return fac.createDoctron(config)
	case DoctronHtml2Svg:
		fac := &(html2SvgFactory{})
		return fac.createDoctron(config)
	case DoctronHtml2Image:
		fac := &(html2ImageFactory{})
		return fac.createDoctron(config)
	case DoctronPdf2Image:
		fac := &(pdf2ImageFactory{})
		return fac.createDoctron(config)
	case DoctronPdfWatermark:
		fac := &(pdfWatermarkFactory{})
		return fac.createDoctron(config)
	default:
		return nil
	}
}

type DoctronFactory interface {
	createDoctron(ctx context.Context, cc converter.ConvertConfig) DoctronI
}

type html2PdfFactory struct {
}

func createHtml2PdfDoctron(config converter.DoctronConfig) *html2pdf {
	return &html2pdf{
		Doctron: Doctron{
			config: config,
		},
	}
}

func (ins *html2PdfFactory) createDoctron(config converter.DoctronConfig) DoctronI {
	return createHtml2PdfDoctron(config)
}

type html2SvgFactory struct {
}

func (ins *html2SvgFactory) createDoctron(config converter.DoctronConfig) DoctronI {
	return &html2svg{
		html2pdf: createHtml2PdfDoctron(config),
	}
}

type html2ImageFactory struct {
}

func (ins *html2ImageFactory) createDoctron(config converter.DoctronConfig) DoctronI {
	return &html2image{
		Doctron: Doctron{
			config: config,
		},
	}
}

type pdf2ImageFactory struct {
}

func (ins *pdf2ImageFactory) createDoctron(config converter.DoctronConfig) DoctronI {
	return &pdf2Image{
		Doctron: Doctron{
			config: config,
		},
	}
}

type pdfWatermarkFactory struct {
}

func (ins *pdfWatermarkFactory) createDoctron(config converter.DoctronConfig) DoctronI {
	return &pdfWatermark{
		Doctron: Doctron{
			config: config,
		},
	}
}
