package transmitsms

import (
	"fmt"
	"github.com/gorilla/schema"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SendSMSRequest struct {
	Message        string         `schema:"message"`
	To             []string       `schema:"-"`
	ListId         *int           `schema:"-"`
	From           string         `schema:"from,omitempty"`
	SendAt         *time.Time     `schema:"-"`
	DlrCallback    string         `schema:"dlr_callback,omitempty"`
	ReplyCallback  string         `schema:"reply_callback,omitempty"`
	Validity       *time.Duration `schema:"-"`
	RepliesToEmail string         `schema:"replies_to_email,omitempty"`
	FromShared     bool           `schema:"-"`
	CountryCode    string         `schema:"countrycode,omitempty"`
}

type SendSMSResponse struct {
	MessageId  string  `json:"message_id"`
	Recipients int     `json:"recipients"`
	Cost       float32 `json:"cost"`
}

func (r *SendSMSRequest) RequestPath() string {
	return "send-sms.json"
}

func (r *SendSMSRequest) ToValues() (v url.Values, err error) {
	v = make(url.Values)
	enc := schema.NewEncoder()
	err = enc.Encode(r, v)
	if err != nil {
		me, ok := err.(schema.MultiError)
		if ok {
			if len(me) != 0 {
				return v, err
			}
			// else if we're here, it's the lame gorilla MultiError killing us.  gah.
			// continue as if nothing went wrong. >_<
			err = nil
		} else {
			return v, err
		}
	}
	// manual encoding time!
	if r.To == nil {
		return v, ErrMalformedRequest
	}
	v.Set("to", strings.Join(r.To, ","))
	if r.ListId != nil {
		v.Set("list_id", strconv.Itoa(*r.ListId))
	}
	if r.SendAt != nil {
		v.Set("send_at", r.SendAt.UTC().Format(SmsTimestampFormat))
	}
	if r.FromShared {
		v.Set("from_shared", "true")
	}
	return v, nil
}

func (r *SendSMSRequest) DecodeResponse(hresp *http.Response) (resp interface{}, err error) {
	rawbody, err := ioutil.ReadAll(hresp.Body)
	hresp.Body.Close()

	if hresp.StatusCode != 200 {
		e := new(ApiError)
		e.HttpCode = hresp.StatusCode
		parts := strings.SplitN(hresp.Status, " ", 2)
		if len(parts) > 1 {
			e.Message = parts[1]
		} else {
			e.Message = fmt.Sprintf("HTTP Error %d", e.HttpCode)
		}
		e.ResponseBody = string(rawbody)

		return nil, e
	}
	return rawbody, nil
}
