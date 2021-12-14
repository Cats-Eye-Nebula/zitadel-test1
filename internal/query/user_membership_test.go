package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/lib/pq"
)

var (
	membershipsStmt = regexp.QuoteMeta(
		"SELECT memberships.user_id" +
			", memberships.roles" +
			", memberships.creation_date" +
			", memberships.change_date" +
			", memberships.sequence" +
			", memberships.resource_owner" +
			", memberships.org_id" +
			", memberships.iam_id" +
			", memberships.project_id" +
			", memberships.grant_id" +
			", zitadel.projections.users_humans.display_name" +
			", zitadel.projections.users_machines.name" +
			", COUNT(*) OVER ()" +
			" FROM (" +
			"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", members.org_id" +
			", NULL::STRING AS iam_id" +
			", NULL::STRING AS project_id" +
			", NULL::STRING AS grant_id" +
			" FROM zitadel.projections.org_members as members" +
			" UNION ALL " +
			"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", NULL::STRING AS org_id" +
			", members.iam_id" +
			", NULL::STRING AS project_id" +
			", NULL::STRING AS grant_id" +
			" FROM zitadel.projections.iam_members as members" +
			" UNION ALL " +
			"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", NULL::STRING AS org_id" +
			", NULL::STRING AS iam_id" +
			", members.project_id" +
			", NULL::STRING AS grant_id" +
			" FROM zitadel.projections.project_members as members" +
			" UNION ALL " +
			"SELECT members.user_id" +
			", members.roles" +
			", members.creation_date" +
			", members.change_date" +
			", members.sequence" +
			", members.resource_owner" +
			", NULL::STRING AS org_id" +
			", NULL::STRING AS iam_id" +
			", members.project_id" +
			", members.grant_id" +
			" FROM zitadel.projections.project_grant_members as members" +
			") AS memberships" +
			" LEFT JOIN zitadel.projections.users_humans ON memberships.user_id = zitadel.projections.users_humans.user_id" +
			" LEFT JOIN zitadel.projections.users_machines ON memberships.user_id = zitadel.projections.users_machines.user_id")
	membershipCols = []string{
		"user_id",
		"roles",
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"org_id",
		"iam_id",
		"project_id",
		"grant_id",
		"display_name",
		"name",
		"count",
	}
)

func Test_MembershipPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareMembershipsQuery no result",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					nil,
					nil,
				),
			},
			object: &Memberships{Memberships: []*Membership{}},
		},
		{
			name:    "prepareMembershipsQuery one org member human",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							"org-id",
							nil,
							nil,
							nil,
							"display name",
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Org:           &OrgMembership{OrgID: "org-id"},
						DisplayName:   "display name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one org member machine",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							"org-id",
							nil,
							nil,
							nil,
							nil,
							"machine-name",
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Org:           &OrgMembership{OrgID: "org-id"},
						DisplayName:   "machine-name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one iam member human",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							"iam-id",
							nil,
							nil,
							"display name",
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						IAM:           &IAMMembership{IAMID: "iam-id"},
						DisplayName:   "display name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one iam member machine",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							"iam-id",
							nil,
							nil,
							nil,
							"machine-name",
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						IAM:           &IAMMembership{IAMID: "iam-id"},
						DisplayName:   "machine-name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one project member human",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							nil,
							"display name",
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Project:       &ProjectMembership{ProjectID: "project-id"},
						DisplayName:   "display name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one project member machine",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							nil,
							nil,
							"machine-name",
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Project:       &ProjectMembership{ProjectID: "project-id"},
						DisplayName:   "machine-name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one project grant member human",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							"grant-id",
							"display name",
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						ProjectGrant: &ProjectGrantMembership{
							GrantID:   "grant-id",
							ProjectID: "project-id",
						},
						DisplayName: "display name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one project grant member machine",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							"grant-id",
							nil,
							"machine-name",
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						ProjectGrant: &ProjectGrantMembership{
							GrantID:   "grant-id",
							ProjectID: "project-id",
						},
						DisplayName: "machine-name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery one for each member type",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueries(
					membershipsStmt,
					membershipCols,
					[][]driver.Value{
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							"org-id",
							nil,
							nil,
							nil,
							"display name",
							nil,
						},
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							"iam-id",
							nil,
							nil,
							"display name",
							nil,
						},
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							nil,
							"display name",
							nil,
						},
						{
							"user-id",
							pq.StringArray{"role1", "role2"},
							testNow,
							testNow,
							uint64(20211202),
							"ro",
							nil,
							nil,
							"project-id",
							"grant-id",
							"display name",
							nil,
						},
					},
				),
			},
			object: &Memberships{
				SearchResponse: SearchResponse{
					Count: 4,
				},
				Memberships: []*Membership{
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Org:           &OrgMembership{OrgID: "org-id"},
						DisplayName:   "display name",
					},
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						IAM:           &IAMMembership{IAMID: "iam-id"},
						DisplayName:   "display name",
					},
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						Project:       &ProjectMembership{ProjectID: "project-id"},
						DisplayName:   "display name",
					},
					{
						UserID:        "user-id",
						Roles:         []string{"role1", "role2"},
						CreationDate:  testNow,
						ChangeDate:    testNow,
						Sequence:      20211202,
						ResourceOwner: "ro",
						ProjectGrant: &ProjectGrantMembership{
							ProjectID: "project-id",
							GrantID:   "grant-id",
						},
						DisplayName: "display name",
					},
				},
			},
		},
		{
			name:    "prepareMembershipsQuery sql err",
			prepare: prepareMembershipsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					membershipsStmt,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
