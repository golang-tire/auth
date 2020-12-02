import React from 'react';
import {Switch, Redirect} from 'react-router-dom';
import {RouteWithLayout, ProtectedRouteWithLayout} from './components';
import {AdminLayout, AuthLayout} from './layouts'
import routes from "./settings";

const Routes = () => {
    return (
        <Switch>
            <Redirect
                exact
                from="/"
                to="/dashboard"
            />
            {
                routes.routeList.map(item => {
                    if (item.protected) {
                        if (item.subMenus) {
                            return (
                                item.subMenus.map(subMenu =>{
                                    console.log("sub", subMenu)
                                    return (
                                        <ProtectedRouteWithLayout
                                            component={subMenu.component}
                                            exact
                                            layout={AdminLayout}
                                            path={subMenu.path}
                                        />
                                    )
                                })
                            )

                        }else {
                            return (
                                <ProtectedRouteWithLayout
                                    component={item.component}
                                    exact
                                    layout={AdminLayout}
                                    path={item.path}
                                />
                            )
                        }
                    } else {
                        return (
                            <RouteWithLayout
                                component={item.component}
                                exact
                                layout={AuthLayout}
                                path="/login"
                            />
                        )
                    }
                })
            }
            <Redirect to="/not-found"/>
        </Switch>
    );
};

export default Routes;
