import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import PropTypes from 'prop-types';
import TokenService from "services/Auth/tokenService";

const ProtectedRouteWithLayout = props => {
    const { layout: Layout, component: Component, ...rest } = props;

    return (
        <Route
            {...rest}
            render={matchProps => (
                TokenService.isAuthenticated() ? (
                    <Layout>
                        <Component {...matchProps} />
                    </Layout>
                ) : (
                    <Redirect
                        to={{
                            pathname: "/login",
                            state: { from: matchProps.location }
                        }}
                    />
                )
            )}
        />
    );
};

ProtectedRouteWithLayout.propTypes = {
    component: PropTypes.any.isRequired,
    layout: PropTypes.any.isRequired,
    path: PropTypes.oneOfType([
        PropTypes.string,
        PropTypes.array
    ])
};

export default ProtectedRouteWithLayout;