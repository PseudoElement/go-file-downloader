package voicechat

func findHost(voiceRoom *VoiceRoom) *User {
	for _, user := range voiceRoom.users {
		if user.isHost {
			return user
		}
	}
	return nil
}
