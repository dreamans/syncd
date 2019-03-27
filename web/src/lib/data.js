import priv from '@/lib/priv'

export default {
    PageSize: 0,
    Page: 0,
    Total: 0,
    DialogSmallWidth: '500px',
    DialogNormalWidth: '750px',
    DialogLargeWidth: '900px',
    DialogNormalTop: '5vh',
    Priv: priv,

    BuildStatusNone: 0, 
    BuildStatusStart: 1,
    BuildStatusSuccess: 2,
    BuildStatusFailed: 3,

    ApplyStatusNone: 1,
    ApplyStatusIng: 2,
    ApplyStatusSuccess: 3,
    ApplyStatusFailed: 4,
    ApplyStatusDrop: 5,
    ApplyStatusRollback: 6,

    DeployModeBranch: 1,
    DeployModelTag: 2,

    BuildStatusNone: 0,
    BuildStatusStart: 1,
    BuildStatusSuccess: 2,
    BuildStatusFailed: 3,

    DeployGroupStatusNone: 0,
    DeployGroupStatusStart: 1,
    DeployGroupStatusSuccess: 2,
    DeployGroupStatusFailed: 3,

    AuditStatusPending: 1,
    AuditStatusOk: 2,
    AuditStatusRefuse: 3,
}