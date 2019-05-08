package pgproto

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
	"github.com/readystock/pgx/chunkreader"
)

type Frontend struct {
	cr *chunkreader.ChunkReader
	w  io.Writer

	// Backend message flyweights
	authentication       Authentication
	backendKeyData       BackendKeyData
	bindComplete         BindComplete
	closeComplete        CloseComplete
	commandComplete      CommandComplete
	copyBothResponse     CopyBothResponse
	copyData             CopyData
	copyInResponse       CopyInResponse
	copyOutResponse      CopyOutResponse
	dataRow              DataRow
	emptyQueryResponse   EmptyQueryResponse
	errorResponse        ErrorResponse
	functionCallResponse FunctionCallResponse
	noData               NoData
	noticeResponse       NoticeResponse
	notificationResponse NotificationResponse
	parameterDescription ParameterDescription
	parameterStatus      ParameterStatus
	parseComplete        ParseComplete
	readyForQuery        ReadyForQuery
	rowDescription       RowDescription

	bodyLen    int
	msgType    byte
	partialMsg bool
}

func NewFrontend(r io.Reader, w io.Writer) (*Frontend, error) {
	cr := chunkreader.NewChunkReader(r)
	return &Frontend{cr: cr, w: w}, nil
}

func (b *Frontend) Send(msg FrontendMessage) error {
	_, err := b.w.Write(msg.Encode(nil))
	return err
}

func (b *Frontend) Receive() (BackendMessage, error) {
	if !b.partialMsg {
		header, err := b.cr.Next(5)
		if err != nil {
			return nil, err
		}

		b.msgType = header[0]
		b.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		b.partialMsg = true
	}

	var msg BackendMessage
	switch b.msgType {
	case PgParseComplete:
		msg = &b.parseComplete
	case PgBindComplete:
		msg = &b.bindComplete
	case PgCloseComplete:
		msg = &b.closeComplete
	case PgNotificationResponse:
		msg = &b.notificationResponse
	case PgCommandComplete:
		msg = &b.commandComplete
	case PgCopyData:
		msg = &b.copyData
	case PgDataRow:
		msg = &b.dataRow
	case PgErrorResponse:
		msg = &b.errorResponse
	case PgCopyInResponse:
		msg = &b.copyInResponse
	case PgCopyOutResponse:
		msg = &b.copyOutResponse
	case PgEmptyQueryResponse:
		msg = &b.emptyQueryResponse
	case PgBackendKeyData:
		msg = &b.backendKeyData
	case PgNoData:
		msg = &b.noData
	case PgNoticeResponse:
		msg = &b.noticeResponse
	case PgAuthentication:
		msg = &b.authentication
	case PgParameterStatus:
		msg = &b.parameterStatus
	case PgParameterDescription:
		msg = &b.parameterDescription
	case PgRowDescription:
		msg = &b.rowDescription
	case PgFunctionCallResponse:
		msg = &b.functionCallResponse
	case PgCopyBothResponse:
		msg = &b.copyBothResponse
	case PgReadyForQuery:
		msg = &b.readyForQuery
	default:
		return nil, errors.Errorf("unknown message type: %c", b.msgType)
	}

	msgBody, err := b.cr.Next(b.bodyLen)
	if err != nil {
		return nil, err
	}

	b.partialMsg = false

	err = msg.Decode(msgBody)
	return msg, err
}
