// internal/protocol/packet_ids.go
// Все константы ID пакетов для протокола 340 (Minecraft 1.12.2).

package protocol

// Clientbound (Сервер -> Клиент)
const (
	// Login state
	SLoginDisconnect   = 0x00
	SEncryptionRequest = 0x01
	SLoginSuccess      = 0x02
	SSetCompression    = 0x03

	// Play state
	SSpawnEntity        = 0x00
	SSpawnExperienceOrb = 0x01
	SSpawnLivingEntity  = 0x03
	SSpawnPlayer        = 0x05
	SEntityAnimation    = 0x06
	SBlockChange        = 0x0B
	SKeepAlive          = 0x1F
	SChunkData          = 0x20
	SJoinGame           = 0x23
	SPlayerPositionLook = 0x2F
	SPlayerAbilities    = 0x2C
	SPlayerListItem     = 0x2E
	SChatMessage        = 0x0F
	SUpdateHealth       = 0x41
	SSpawnPosition      = 0x46
	SSetSlot            = 0x16
	SWindowItems        = 0x14
	SDisconnect         = 0x1A
	SEntityMetadata     = 0x3C
	SEntityVelocity     = 0x3E
	SEntityTeleport     = 0x4C
	STimeUpdate         = 0x47
	SRespawn            = 0x35
)

// Serverbound (Клиент -> Сервер)
const (
	// Login state
	CLoginStart         = 0x00
	CEncryptionResponse = 0x01

	// Play state
	CTeleportConfirm      = 0x00
	CChatMessage          = 0x02
	CClientStatus         = 0x03
	CClientSettings       = 0x04
	CKeepAlive            = 0x0B
	CPlayerPosition       = 0x0D
	CPlayerPositionLook   = 0x0E
	CPlayerLook           = 0x0F
	CPlayerMovement       = 0x0C
	CPlayerDigging        = 0x14
	CEntityAction         = 0x15
	CPlayerBlockPlacement = 0x1F
	CUseItem              = 0x20
	CHeldItemChange       = 0x1A
	CClickWindow          = 0x08
	CCloseWindow          = 0x09
	CAnimation            = 0x1D
)