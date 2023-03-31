package xtextract

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
	"github.com/grokify/awsgo/xtextract/ocrutil"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/image/imageutil"
)

/*

AWS: https://github.com/aws-samples/amazon-textract-code-samples
GOOGLE: https://cloud.google.com/document-ai/docs/process-documents-client-libraries

	//	*ProcessRequest_InlineDocument
	//	*ProcessRequest_RawDocument

*/

type TextractClientMore struct {
	textractClient *textract.Textract
}

func NewTextractClientMore(sess *session.Session, region string) (*TextractClientMore, error) {
	c, err := NewClient(sess, region)
	if err != nil {
		return nil, err
	}
	return &TextractClientMore{textractClient: c}, nil
}

func (cm TextractClientMore) DetectText(ctx context.Context, b []byte, opts ...request.Option) (*textract.DetectDocumentTextOutput, error) {
	if len(b) == 0 {
		return nil, ocrutil.ErrNoInputData
	}
	return cm.DetectTextFromDocument(ctx, DocumentFromBytes(b), opts...)
}

func (cm TextractClientMore) DetectTextFromFile(ctx context.Context, filename string, opts ...request.Option) (*textract.DetectDocumentTextOutput, error) {
	doc, err := DocumentFromFilename(filename)
	if err != nil {
		return nil, err
	}
	return cm.DetectTextFromDocument(ctx, doc, opts...)
}

func (cm TextractClientMore) DetectTextFromImage(ctx context.Context, img image.Image, opts ...request.Option) (*textract.DetectDocumentTextOutput, error) {
	doc, err := DocumentFromImage(img)
	if err != nil {
		return nil, err
	}
	return cm.DetectTextFromDocument(ctx, doc, opts...)
}

func (cm TextractClientMore) DetectTextFromImageLocationsVertical(ctx context.Context, imglocations []string, opts ...request.Option) (*textract.DetectDocumentTextOutput, error) {
	img, err := imageutil.MergeYSameXRead(imglocations, true)
	if err != nil {
		return nil, err
	}
	return cm.DetectTextFromImage(ctx, img, opts...)
}

func (cm TextractClientMore) OCRSync(ctx context.Context, b []byte) (ocrutil.TextResults, error) {
	return cm.outputToTextResults(cm.DetectText(ctx, b))
}

func (cm TextractClientMore) OCRSyncImageWriteFile(ctx context.Context, outfilename string, perm os.FileMode, imgs []image.Image) (ocrutil.TextResults, error) {
	if len(imgs) == 0 {
		return ocrutil.TextResults{}, errors.New("no images")
	}
	img := imgs[0]
	if len(imgs) > 1 {
		img = imageutil.MergeYSameX(imgs, true)
	}
	tr, err := cm.outputToTextResults(cm.DetectTextFromImage(ctx, img))
	if err != nil {
		return tr, err
	}
	b, err := jsonutil.MarshalSimple(tr, "", "  ")
	if err != nil {
		return tr, err
	}
	err = os.WriteFile(outfilename, b, perm)
	return tr, err
}

func (cm TextractClientMore) OCRSyncImageLocationsWriteFile(ctx context.Context, outfilename string, perm os.FileMode, imglocations []string) (ocrutil.TextResults, error) {
	tr, err := cm.outputToTextResults(cm.DetectTextFromImageLocationsVertical(ctx, imglocations))
	if err != nil {
		return tr, err
	}
	b, err := jsonutil.MarshalSimple(tr, "", "  ")
	if err != nil {
		return tr, err
	}
	err = os.WriteFile(outfilename, b, perm)
	return tr, err
}

func (cm TextractClientMore) outputToTextResults(output *textract.DetectDocumentTextOutput, err error) (ocrutil.TextResults, error) {
	if err != nil {
		return ocrutil.TextResults{}, err
	}
	blks := Blocks(output.Blocks)
	tr := blks.TextResults()
	tr.OCRDateTime = time.Now().UTC()
	tr.OCRService = "aws"
	return tr, nil
}

func (cm TextractClientMore) DetectTextFromDocument(ctx context.Context, doc *textract.Document, opts ...request.Option) (*textract.DetectDocumentTextOutput, error) {
	// See: https://docs.aws.amazon.com/sdk-for-go/api/service/textract/#Textract.DetectDocumentTextRequest
	input := &textract.DetectDocumentTextInput{}
	input = input.SetDocument(doc)
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	if ctx == nil && len(opts) == 0 {
		req, resp := cm.textractClient.DetectDocumentTextRequest(input)
		// (req *request.Request, output *DetectDocumentTextOutput)

		err = req.Send()
		if err != nil {
			return resp, err
		}
		// resp is now filled
		fmt.Println(resp)
		return resp, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return cm.textractClient.DetectDocumentTextWithContext(ctx, input, opts...)

	// GetDocumentTextDetectionOutput - https://docs.aws.amazon.com/sdk-for-go/api/service/textract/#GetDocumentTextDetectionOutput
}

// https://docs.aws.amazon.com/textract/latest/dg/how-it-works-detecting.html

func NewClient(sess *session.Session, region string) (*textract.Textract, error) {
	region = strings.TrimSpace(region)
	if region == "" {
		return nil, errors.New("region must be supplied")
	}
	var err error
	if sess == nil {
		sess, err = session.NewSession()
		if err != nil {
			return nil, err
		}
	}
	// mySession := session.Must(session.NewSession())
	region = strings.TrimSpace(region)
	if region == "" {
		return textract.New(sess), nil
	}
	return textract.New(sess, aws.NewConfig().WithRegion(region)), nil
}

func DocumentFromFilename(name string) (*textract.Document, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return DocumentFromBytes(b), nil
}

func DocumentFromImage(i image.Image) (*textract.Document, error) {
	b, err := imageutil.BytesJPEG(i, nil)
	if err != nil {
		return nil, err
	}
	return DocumentFromBytes(b), nil
}

// DocumentFromBytes - The document bytes must be in PNG or JPEG format.
func DocumentFromBytes(b []byte) *textract.Document {
	return &textract.Document{
		Bytes: b,
	}
}

func ReadFileDetectDocumentTextOutput(filename string) (*textract.DetectDocumentTextOutput, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	output := &textract.DetectDocumentTextOutput{}
	err = json.Unmarshal(b, output)
	return output, err
}

/*
func WriteFileStringer(name string, s fmt.Stringer, perm os.FileMode) error {
	return os.WriteFile(name, []byte(s.String()), perm)
}
*/

// *textract.DetectDocumentTextOutput
