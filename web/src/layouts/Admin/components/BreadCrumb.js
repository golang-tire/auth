import React from "react";
import { Link, withRouter } from 'react-router-dom';
import {Breadcrumb} from 'antd';
import settings from "settings";

import {DashboardOutlined} from '@ant-design/icons';

const breadcrumbNameMap = {}
settings.routeList.map(item =>{
    if (item.subRoutes){
        item.subRoutes.map(s =>{
            breadcrumbNameMap[s.path] = s.name
        })
    }else{
        breadcrumbNameMap[item.path] = item.name
    }
})

const BreadCrumb = withRouter(props => {
    const { location } = props;
    const pathSnippets = location.pathname.split('/').filter(i => i);
    const extraBreadcrumbItems = pathSnippets.map((_, index) => {
        const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;
        return (
            <Breadcrumb.Item key={url}>
                <Link to={url}>{breadcrumbNameMap[url]}</Link>
            </Breadcrumb.Item>
        );
    });
    const breadcrumbItems = [
        <Breadcrumb.Item key="dashboard">
            <Link to="/dashboard"><DashboardOutlined /></Link>
        </Breadcrumb.Item>,
    ].concat(extraBreadcrumbItems);

    return(
        <Breadcrumb>{breadcrumbItems}</Breadcrumb>
   )
})

export default BreadCrumb;