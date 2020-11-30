import React from "react";
import { Result, Button } from 'antd';
import {Link} from "react-router-dom";

const NotFound = props => {

    return (
        <Result
            status="404"
            title="404"
            subTitle="Sorry, the page you visited does not exist."
            extra={<Button type="primary"><Link to="/dashboard">go to Dashboard</Link></Button>}
        />
    )
}

export default NotFound;