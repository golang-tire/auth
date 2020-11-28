import React from 'react';
import {Switch, Redirect} from 'react-router-dom';
import {RouteWithLayout, ProtectedRouteWithLayout} from './components';
import {AdminLayout, AuthLayout} from './layouts'
import {
    Dashboard,
    Users,
    UserEdit,
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
                path={["/users/edit", "/users/edit/:userUuid"]}
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
