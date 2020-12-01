import {
    Dashboard,
    Users,
    UserEdit,
    Domains,
    DomainEdit,
    Roles,
    RoleEdit,
    Rules,
    RuleEdit,
    AuditLogs,
    Login,
    NotFound,
} from './pages'

export default {
    routeList: [
        {
            id: '1',
            icon: 'dashboard',
            name: 'Dashboard',
            path: '/dashboard',
            component: Dashboard,
            protected: true,
            sideMenu: true
        },
        {
            id: '2',
            icon: 'users',
            name: 'Users',
            path: '/users',
            component: Users,
            protected: true,
            sideMenu: true
        },
        {
            path: ["/users/edit", "/users/edit/:Uuid"],
            component: UserEdit,
            protected: true,
            sideMenu: false
        },
        {
            id: '3',
            icon: 'domains',
            name: 'Domains',
            path: '/domains',
            component: Domains,
            protected: true,
            sideMenu: true
        },
        {
            path: ["/domains/edit", "/domains/edit/:Uuid"],
            component: DomainEdit,
            protected: true,
            sideMenu: false
        },
        {
            id: '4',
            icon: 'roles',
            name: 'Roles',
            path: '/roles',
            component: Roles,
            protected: true,
            sideMenu: true
        },
        {
            path: ["/roles/edit", "/roles/edit/:Uuid"],
            component: RoleEdit,
            protected: true,
            sideMenu: false
        },
        {
            id: '5',
            icon: 'rules',
            name: 'Rules',
            path: '/rules',
            component: Rules,
            protected: true,
            sideMenu: true
        },
        {
            path: ["/rules/edit", "/rules/edit/:Uuid"],
            component: RuleEdit,
            protected: true,
            sideMenu: false
        },
        {
            id: '6',
            icon: 'audit-logs',
            name: 'Audit Logs',
            path: '/audit-logs',
            component: AuditLogs,
            protected: true,
            sideMenu: true
        },
        {
            path: "/login",
            component: Login,
            protected: false,
            sideMenu: false
        },
        {
            path: "/not-found",
            component: NotFound,
            protected: true,
            sideMenu: false
        },
    ],
}