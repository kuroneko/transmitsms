package transmitsms

import (
	"github.com/gorilla/schema"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SendSMSRequest struct {
	Message        string         `schema:"message"`
	To             []string       `schema:"-"`
	ListId         *int           `schema:"list_id,omitempty"`
	From           string         `schema:"from,omitempty"`
	SendAt         *time.Time     `schema:"-"`
	DlrCallback    string         `schema:"dlr_callback,omitempty"`
	ReplyCallback  string         `schema:"reply_callback,omitempty"`
	Validity       *time.Duration `schema:"-"`
	RepliesToEmail string         `schema:"replies_to_email,omitempty"`
	FromShared     bool           `schema:"-"`
	CountryCode    string         `schema:"countrycode,omitempty"`
}

func (r *SendSMSRequest) RequestPath() string {
	return "send-sms.json"
}

func (r *SendSMSRequest) ToValues() (v url.Values, err error) {
	v = make(url.Values)
	enc := schema.NewEncoder()
	err = enc.Encode(r, v)
	if err != nil {
		return v, err
	}
	// manual encoding time!
	if r.To == nil {
		return v, ErrMalformedRequest
	}
	v.Set("to", strings.Join(r.To, ","))
	if r.SendAt != nil {
		v.Set("send_at", r.SendAt.UTC().Format(SmsTimestampFormat))
	}
	if r.FromShared {
		v.Set("from_shared", "true")
	}
	return v, nil
}

func (r *SendSMSRequest) DecodeResponse(hresp *http.Response) (resp interface{}, err error) {
	return nil, nil
}
