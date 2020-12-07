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
        let test = breadcrumbNameMap[url]
        let clickable = true
        if (pathSnippets.length === 2 && url.includes("edit")){
            test = "new"
            clickable = false
        }else if (pathSnippets.length === 3 && url.includes("edit") && index===2) {
            test = pathSnippets[index]
        }else if (pathSnippets.length === 3 && url.includes("edit")) {
            test = "detail"
            clickable = false
        }
        return (
            <Breadcrumb.Item key={url}>
                {clickable&&<Link to={url}>{test}</Link>}
                {!clickable&&<span>{test}</span>}
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