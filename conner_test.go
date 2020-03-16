package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExternalProcedure(t *testing.T) {
	type args struct {
		connectors []Connector
	}
	tests := []struct {
		name string
		args args
		want *ExternalProcedure
	}{
		{
			args: args{
				connectors: []Connector{
					NewMySQLConfig(
						MySQLHost("127.0.0.1"),
						MySQLPort("3306"),
						MySQLUsername("root"),
						MySQLPassword(""),
						MySQLDatabase("test"),
					),
				},
			},
			want: &ExternalProcedure{Connectors: []Connector{
				&MySQLConfig{
					host:     "127.0.0.1",
					port:     "3306",
					user:     "root",
					password: "",
					database: "test",
				},
			}},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {

			connectors := NewExternalProcedure(testCase.args.connectors...)
			for index, database := range connectors.Connectors {
				mysql, ok := database.(*MySQLConfig)
				assert.True(t, ok)
				want, ok := testCase.want.Connectors[index].(*MySQLConfig)
				assert.True(t, ok)
				assert.Equal(t, want.host, mysql.host)
				assert.Equal(t, want.port, mysql.port)
				assert.Equal(t, want.user, mysql.user)
				assert.Equal(t, want.password, mysql.password)
				assert.Equal(t, want.database, mysql.database)
			}
		})
	}
}
