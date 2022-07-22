package controller

import (
	"context"
	"errors"
	"time"

	uuid "github.com/iris-contrib/go.uuid"
	"github.com/kataras/iris/v12"
	irisContext "github.com/kataras/iris/v12/context"
	"github.com/lampnick/doctron/common"
	"github.com/lampnick/doctron/conf"
	"github.com/lampnick/doctron/converter/doctron_core"
	"github.com/lampnick/doctron/worker"
)

func Html2SvgHandler(ctx iris.Context) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(conf.LoadedConfig.Doctron.ConvertTimeout)*time.Second)
	defer cancel()
	traceId, _ := uuid.NewV4()
	outputDTO := common.NewDefaultOutputDTO(nil)
	doctronConfig, err := initHtml2PdfConfig(ctx)
	if err != nil {
		outputDTO.Code = common.InvalidParams
		outputDTO.Message = err.Error()
		_, _ = common.NewJsonOutput(ctx, outputDTO)
		return
	}
	doctronConfig.TraceId = traceId
	doctronConfig.Ctx = ctxTimeout
	doctronConfig.IrisCtx = ctx
	doctronConfig.DoctronType = doctron_core.DoctronHtml2Svg

	doctronOutputDTO, err := worker.Pool.ProcessTimed(doctronConfig, time.Duration(conf.LoadedConfig.Doctron.ConvertTimeout)*time.Second)
	if err != nil {
		outputDTO.Code = common.ConvertSvgFailed
		outputDTO.Message = "worker run process failed." + err.Error()
		_, _ = common.NewJsonOutput(ctx, outputDTO)
		return
	}
	doctronOutput, ok := doctronOutputDTO.(worker.DoctronOutputDTO)
	if !ok {
		outputDTO.Code = common.ConvertSvgFailed
		outputDTO.Message = "error type assert to DoctronOutputDTO"
		_, _ = common.NewJsonOutput(ctx, outputDTO)
		return
	}

	if errors.Is(doctronOutput.Err, worker.ErrNoNeedToUpload) {
		ctx.Header(irisContext.ContentTypeHeaderKey, "image/svg+xml")
		_, err = ctx.Write(doctronOutput.Buf)
		if err != nil {
			outputDTO.Code = common.ConvertSvgFailed
			_, _ = common.NewJsonOutput(ctx, outputDTO)
			return
		}
		return
	}
	if doctronOutput.Err != nil {
		outputDTO.Code = common.ConvertSvgFailed
		outputDTO.Message = doctronOutput.Err.Error()
		_, _ = common.NewJsonOutput(ctx, outputDTO)
		return
	}
	outputDTO.Data = doctronOutput.Url
	_, _ = common.NewJsonOutput(ctx, outputDTO)
	return
}
