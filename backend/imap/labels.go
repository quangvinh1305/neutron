package imap

import (
	"errors"

	"github.com/emersion/neutron/backend"
)

func getLabelID(mailbox string) string {
	lbl := mailbox
	switch lbl {
	case "INBOX":
		lbl = backend.InboxLabel
	case "Draft", "Drafts":
		lbl = backend.DraftLabel
	case "Sent":
		lbl = backend.SentLabel
	case "Trash":
		lbl = backend.TrashLabel
	case "Spam", "Junk":
		lbl = backend.SpamLabel
	case "Archive", "Archives":
		lbl = backend.ArchiveLabel
	case "Important", "Starred":
		lbl = backend.StarredLabel
	}
	return lbl
}

var colors = []string{
	// Dark
	"#7272a7",
	"#cf5858",
	"#c26cc7",
	"#7569d1",
	"#69a9d1",
	"#5ec7b7",
	"#72bb75",
	"#c3d261",
	"#e6c04c",
	"#e6984c",

	// Light
	"#8989ac",
	"#cf7e7e",
	"#c793ca",
	"#9b94d1",
	"#a8c4d5",
	"#97c9c1",
	"#9db99f",
	"#c6cd97",
	"#e7d292",
	"#dfb28",
}

func getLabelColor(i int) string {
	return colors[i % len(colors)]
}

type Labels struct {
	*conns
}

func (b *Labels) ListLabels(user string) (labels []*backend.Label, err error) {
	mailboxes, err := b.getMailboxes(user)
	if err != nil {
		return
	}

	i := 0
	for _, mailbox := range mailboxes {
		name := mailbox.Name

		if getLabelID(name) != name {
			continue // This is a system mailbox, not a custom one
		}

		color := getLabelColor(i)

		labels = append(labels, &backend.Label{
			ID: name,
			Name: name,
			Color: color,
			Display: 1,
			Order: i,
		})

		i++
	}

	return
}

func (b *Labels) InsertLabel(user string, label *backend.Label) (inserted *backend.Label, err error) {
	labels, err := b.ListLabels(user)
	if err != nil {
		return
	}

	c, unlock, err := b.getConn(user)
	if err != nil {
		return
	}
	defer unlock()

	_, _, err = wait(c.Create(label.Name))
	if err != nil {
		return
	}

	i := len(labels)
	inserted = label
	label.ID = label.Name
	label.Color = getLabelColor(i)
	label.Order = i
	return
}

func (b *Labels) UpdateLabel(user string, update *backend.LabelUpdate) (*backend.Label, error) {
	return nil, errors.New("Not yet implemented")
}

func (b *Labels) DeleteLabel(user, id string) error {
	return errors.New("Not yet implemented")
}

func newLabels(conns *conns) *Labels {
	return &Labels{
		conns: conns,
	}
}
