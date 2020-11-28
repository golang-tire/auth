import React from 'react';
import {Switch, Redirect} from 'react-router-dom';
import {RouteWithLayout, ProtectedRouteWithLayout} from './components';
import {AdminLayout, AuthLayout} from './layouts'
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
    Login,
    NotFound,
} from './pages'

const Routes = () => {
    return (
        <Switch>
            <Redirect
                exact
                from="/"
                to="/dashboard"
            />
            <ProtectedRouteWithLayout
                component={Dashboard}
                exact
                layout={AdminLayout}
                path="/dashboard"
            />
            <ProtectedRouteWithLayout
                component={Users}
                exact
                layout={AdminLayout}
                path="/users"
            />
            <ProtectedRouteWithLayout
                component={UserEdit}
                exact
                layout={AdminLayout}
                path={["/users/edit", "/users/edit/:Uuid"]}
            />
            <ProtectedRouteWithLayout
                component={Domains}
                exact
                layout={AdminLayout}
                path="/domains"
            />
            <ProtectedRouteWithLayout
                component={DomainEdit}
                exact
                layout={AdminLayout}
                path={["/domains/edit", "/domains/edit/:Uuid"]}
            />
            <ProtectedRouteWithLayout
                component={Roles}
                exact
                layout={AdminLayout}
                path="/roles"
            />
            <ProtectedRouteWithLayout
                component={RoleEdit}
                exact
                layout={AdminLayout}
                path={["/roles/edit", "/roles/edit/:Uuid"]}
            />
            <ProtectedRouteWithLayout
                component={Rules}
                exact
                layout={AdminLayout}
                path="/rules"
            />
            <ProtectedRouteWithLayout
                component={RuleEdit}
                exact
                layout={AdminLayout}
                path={["/rules/edit", "/rules/edit/:Uuid"]}
            />
            <RouteWithLayout
                component={Login}
                exact
                layout={AuthLayout}
                path="/login"
            />
            <ProtectedRouteWithLayout
                component={NotFound}
                exact
                layout={AdminLayout}
                path="/not-found"
            />
            <Redirect to="/not-found"/>
        </Switch>
    );
};

export default Routes;
