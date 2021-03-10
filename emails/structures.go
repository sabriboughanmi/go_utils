package emails


import "net/textproto"


//EmailAddress contains all data related to an Email address for Server Side Usage
type EmailAddress struct {
	address  string
	userName string
	port     int
	password string
	host     string
}


// Email is the type used for email messages
type Email struct {
	ReplyTo     []string
	To          []string
	Bcc         []string
	Cc          []string
	Subject     string
	Text        []byte // Plaintext message (optional)
	HTML        []byte // Html message (optional)
	Sender      string // override From as SMTP envelope sender (optional)
	Headers     textproto.MIMEHeader
	Attachments []*Attachment
	ReadReceipt []string
}

// part is a copyable representation of a multipart.Part
type part struct {
	header textproto.MIMEHeader
	body   []byte
}

// Attachment is a struct representing an email attachment.
// Based on the mime/multipart.FileHeader struct, Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment struct {
	Filename    string
	ContentType string
	Header      textproto.MIMEHeader
	Content     []byte
	HTMLRelated bool
}

