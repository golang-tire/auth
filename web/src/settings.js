import {
    Dashboard,
    Users,
    UserEdit,
    Domains,
    DomainEdit,
    Apps,
    AppEdit,
    Resources,
    ResourceEdit,
    Objects,
    ObjectEdit,
    Roles,
    RoleEdit,
    Rules,
    RuleEdit,
    AuditLogs,
    Login,
    NotFound,
} from './pages'

export default {
    defaultSelectedMenus: ["dashboard"],
    routeList: [
        {
            id: 'dashboard',
            icon: 'dashboard',
            name: 'Dashboard',
            path: '/dashboard',
            component: Dashboard,
            protected: true,
            sideMenu: true
        },
        {
            id: 'workspaces',
            icon: 'domains',
            name: 'Workspaces',
            protected: true,
            sideMenu: true,
            subRoutes: [
                {
                    id: 'domains',
                    name: 'Domains',
                    path: '/domains',
                    component: Domains,
                    sideMenu: true,
                },
                {
                    id:"domains-edit",
                    path: ["/domains/edit", "/domains/edit/:Uuid"],
                    component: DomainEdit,
                    sideMenu: false
                },
            ]
        },
        {
            id: 'user-management',
            icon: 'users',
            name: 'User Management',
            protected: true,
            sideMenu: true,
            subRoutes: [
                {
                    id: 'users',
                    name: 'Users',
                    path: '/users',
                    component: Users,
                    sideMenu: true,
                },
                {
                    id: "users-edit",
                    path: ["/users/edit", "/users/edit/:Uuid"],
                    component: UserEdit,
                    sideMenu: false
                },
                {
                    id: 'roles',
                    icon: 'roles',
                    name: 'Roles',
                    path: '/roles',
                    component: Roles,
                    sideMenu: true
                },
                {
                    id: "roles-edit",
                    name: 'Roles',
                    path: ["/roles/edit", "/roles/edit/:Uuid"],
                    component: RoleEdit,
                    sideMenu: false
                },
                {
                    id: 'rules',
                    icon: 'rules',
                    name: 'Rules',
                    path: '/rules',
                    component: Rules,
                    sideMenu: true
                },
                {
                    id: "rules-edit",
                    path: ["/rules/edit", "/rules/edit/:Uuid"],
                    component: RuleEdit,
                    sideMenu: false
                },
            ]
        },
        {
            id: 'applications',
            icon: 'apps',
            name: 'Applications',
            protected: true,
            sideMenu: true,
            subRoutes:[
                {
                    id: 'apps',
                    name: 'Apps',
                    path: '/apps',
                    component: Apps,
                    sideMenu: true
                },
                {
                    id: "apps-edit",
                    path: ["/apps/edit", "/apps/edit/:Uuid"],
                    component: AppEdit,
                    sideMenu: false
                },
                {
                    id: 'resources',
                    name: 'Resources',
                    path: '/resources',
                    component: Resources,
                    sideMenu: true
                },
                {
                    id: "resources-edit",
                    path: ["/resources/edit", "/resources/edit/:Uuid"],
                    component: ResourceEdit,
                    sideMenu: false
                },
                {

                    id: 'objects',
                    name: 'Objects',
                    path: '/objects',
                    component: Objects,
                    sideMenu: true
                },
                {
                    id: "objects-edit",
                    path: ["/objects/edit", "/objects/edit/:Uuid"],
                    component: ObjectEdit,
                    sideMenu: false
                },
            ]
        },
        {
            id: 'audit-logs',
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