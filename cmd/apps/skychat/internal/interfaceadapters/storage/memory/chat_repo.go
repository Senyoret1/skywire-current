package memory

import (
	"fmt"
	"sync"

	"github.com/skycoin/skywire-utilities/pkg/cipher"
	"github.com/skycoin/skywire/cmd/apps/skychat/internal/domain/chat"
)

//ChatRepo Implements the Repository Interface to provide an in-memory storage provider
type ChatRepo struct {
	chats   map[cipher.PubKey]chat.Chat
	chatsMu sync.Mutex
}

//NewRepo Constructor
func NewChatRepo() *ChatRepo {
	cR := ChatRepo{}
	cR.chats = make(map[cipher.PubKey]chat.Chat)
	return &cR
}

//GetByPK Returns the chat with the provided pk
func (r *ChatRepo) GetByPK(pk cipher.PubKey) (*chat.Chat, error) {
	r.chatsMu.Lock()
	defer r.chatsMu.Unlock()

	chat, ok := r.chats[pk]
	if !ok {
		return nil, fmt.Errorf("chat not found")
	}
	return &chat, nil
}

//GetAll Returns all stored chats
func (r *ChatRepo) GetAll() ([]chat.Chat, error) {
	r.chatsMu.Lock()
	defer r.chatsMu.Unlock()

	keys := make([]cipher.PubKey, 0)

	for key := range r.chats {
		keys = append(keys, key)
	}

	var values []chat.Chat
	for _, value := range r.chats {
		values = append(values, value)
	}

	return values, nil
}

//Add the provided chat
func (r *ChatRepo) Add(chat chat.Chat) error {
	r.chatsMu.Lock()
	defer r.chatsMu.Unlock()

	r.chats[chat.GetPK()] = chat
	return nil
}

//Update the provided chat
func (r *ChatRepo) Update(chat chat.Chat) error {
	r.chatsMu.Lock()
	defer r.chatsMu.Unlock()

	r.chats[chat.GetPK()] = chat
	return nil
}

//Delete the chat with the provided pk
func (r *ChatRepo) Delete(pk cipher.PubKey) error {
	r.chatsMu.Lock()
	defer r.chatsMu.Unlock()

	_, exists := r.chats[pk]
	if !exists {
		return fmt.Errorf("id %v not found", pk.String())
	}
	delete(r.chats, pk)
	return nil
}

/*
func NewDummyChatRepo(pk cipher.PubKey) *ChatRepo {
	cR := ChatRepo{}

	defaultpng := "iVBORw0KGgoAAAANSUhEUgAAAPoAAAD6CAIAAAAHjs1qAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyZpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuNS1jMDIxIDc5LjE1NTc3MiwgMjAxNC8wMS8xMy0xOTo0NDowMCAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTQgKFdpbmRvd3MpIiB4bXBNTTpJbnN0YW5jZUlEPSJ4bXAuaWlkOjY5OTMyNUUzOTY0QjExRUI4MDZERkQ5M0JBOUY1NThGIiB4bXBNTTpEb2N1bWVudElEPSJ4bXAuZGlkOjY5OTMyNUU0OTY0QjExRUI4MDZERkQ5M0JBOUY1NThGIj4gPHhtcE1NOkRlcml2ZWRGcm9tIHN0UmVmOmluc3RhbmNlSUQ9InhtcC5paWQ6Njk5MzI1RTE5NjRCMTFFQjgwNkRGRDkzQkE5RjU1OEYiIHN0UmVmOmRvY3VtZW50SUQ9InhtcC5kaWQ6Njk5MzI1RTI5NjRCMTFFQjgwNkRGRDkzQkE5RjU1OEYiLz4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz7FIRCzAAAPn0lEQVR42uyd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"

	//Origin1
	pk1 := cipher.PubKey{}
	pk1.Set("0200369389087217c9e9b5936c329c120c7ed1b29b404d708bc88634e4d6b4966b")
	i1 := info.NewInfo(pk1, "Peer1", "Peer1 Description", defaultpng)
	//TODO: make chat.NewChat()
	//make dummy chats
	PC1msg1 := message.NewTextMessage(pk1, []byte("Hello"))
	PC1msg2 := message.NewTextMessage(pk, []byte("Hello Back"))
	PC1 := chat.NewChat(pk1, chat.PeerChat, i1, []message.Message{PC1msg1, PC1msg2})

	//Origin2
	pk2 := cipher.PubKey{}
	pk2.Set("0200369389087217c9e9b5936c329c120c7ed1b29b404d708bc88634e4d6b4966a")
	i2 := info.NewInfo(pk2, "Peer2", "Peer2 Description", defaultpng)
	//TODO: make chat.NewChat()
	//make dummy chats
	PC2msg1 := message.NewTextMessage(pk2, []byte("Hello Peer2"))
	PC2msg2 := message.NewTextMessage(pk, []byte("Hello Back Peer2"))
	PC2 := chat.NewChat(pk2, chat.PeerChat, i2, []message.Message{PC2msg1, PC2msg2})

	cR.chats = make(map[cipher.PubKey]chat.Chat)
	cR.chats[pk1] = PC1
	cR.chats[pk2] = PC2

	return &cR
}
*/
