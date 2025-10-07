package v2

import (
	_ "time"

	v120 "github.com/sealdice-ce/sealdice-ce/core/migrate/v2/v120"
	v131 "github.com/sealdice-ce/sealdice-ce/core/migrate/v2/v131"
	v141 "github.com/sealdice-ce/sealdice-ce/core/migrate/v2/v141"
	v144 "github.com/sealdice-ce/sealdice-ce/core/migrate/v2/v144"
	v150 "github.com/sealdice-ce/sealdice-ce/core/migrate/v2/v150"
	operator "github.com/sealdice-ce/sealdice-ce/core/utils/dboperator/engine"
	upgrade "github.com/sealdice-ce/sealdice-ce/core/utils/upgrader"
	"github.com/sealdice-ce/sealdice-ce/core/utils/upgrader/store"
)

func InitUpgrader(operator operator.DatabaseOperator) error {
	store := store.NewJSONStore("upgrade_metadata.json")
	mgr := &upgrade.Manager{Store: store, Database: operator}
	// V120注册
	mgr.Register(v120.V120Migration)
	mgr.Register(v120.V120LogMessageMigration)
	// V131注册
	mgr.Register(v131.V131ConfigUpdateMigration)
	// V141注册
	mgr.Register(v141.V141ConfigUpdateMigration)
	// v144注册
	mgr.Register(v144.V144RemoveOldHelpDocMigration)
	// v150注册
	mgr.Register(v150.V150UpgradeAttrsMigration)
	mgr.Register(v150.V150FixGroupInfoMigration)
	err := mgr.ApplyAll()
	if err != nil {
		return err
	}
	return nil
}
