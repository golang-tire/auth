import React from "react";
import PropTypes from "prop-types";

const Auth = props => {
    const {children} = props;

    return (
        <div>
            {children}
        </div>
    )
}

Auth.propTypes = {
    children: PropTypes.node
};

export default Auth;