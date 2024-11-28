package models

type VANT struct {
	ID             int
	X              float64
	Y              float64
	Z              float64
	MessagesBuffer map[string][]int
}

func (v *VANT) ReceiveMessage(messageId string, message []int) {
	if v.MessagesBuffer == nil {
		v.MessagesBuffer = make(map[string][]int)
	}
	if v.MessagesBuffer[messageId] == nil {
		v.MessagesBuffer[messageId] = message
	} else {
		v.MessagesBuffer[messageId] = append(v.MessagesBuffer[messageId], message...)
	}
}

func bitsToString(bits []int) string {
	var result string
	for i := 0; i < len(bits); i += 8 {
		// Pega 8 bits e converte para um byte (caractere ASCII)
		var byteValue byte
		for j := 0; j < 8 && i+j < len(bits); j++ {
			byteValue |= byte(bits[i+j]) << (7 - j)
		}
		result += string(byteValue)
	}
	return result
}

func (v *VANT) GetMessages() []string {
	var messages []string
	for _, message := range v.MessagesBuffer {
		messages = append(messages, bitsToString(message))
	}
	return messages
}

func (v *VANT) HasMessage(message string) bool {
	for id := range v.MessagesBuffer {
		if id == message {
			return true
		}
	}

	return false
}
