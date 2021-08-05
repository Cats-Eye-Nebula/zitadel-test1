---
title: Setup
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

We also provide services to setup our ZITADEL with the operators also provided by us.

<Tabs 
    defaultValue="zitadel"
    values={[
        {label: 'ZITADEL Setup with Kubernetes', value: 'zitadel'}, 
        {label: 'ZITADEL Setup with ORBOS', value: 'orbos'}, 
        {label: 'Checkup', value: 'checkup'},
    ]}>
    <TabItem value="zitadel">
        <p>In Scope</p>
        <ul>
            <li>Check prerequisites and architecture</li>
            <li>Installation and configuration of ZITADEL with the ZITADEL-operator</li>
            <li>Installation and configuration of CockroachDB with the Database-operator</li>
            <li>Functional testing of the ZITADEL instance</li>
        </ul>
        <p>Out of Scope</p>
        <ul>
            <li>Integration into internal monitoring and alerting</li>
            <li>Multi-cluster architecture deployments</li>
            <li>DNS, Network and Firewall configuration</li>
            <li>Kubernetes configuration</li>
            <li>Changes for specific environments</li>
            <li>Performance testing</li>
            <li>Production deployment</li>
            <li>Application-side coding, configuration, or tuning</li>
            <li>Changes or configuration on assets used in ZITADEL</li>
            <li>Setting up or maintaining backup</li>
        </ul>
        <p>Prerequisites</p>
        <ul>
            <li>Running Kubernetes with possibility to deploy to namespaces caos-system and caos-zitadel</li>
            <li>S3-storage for assets in ZITADEL</li>
            <li>S3-storage for backups of the database</li>
            <li>CDN or internal DNS to create DNS records for 4 subdomains pointing at ZITADEL</li>
            <li>Email relay to send emails from ZITADEL</li>
            <li>Account for Twilio to send SMS from ZITADEL</li>
        </ul>
        <p>Deliverable</p>
        <ul>
            <li>Running CockroachDB</li>
            <li>Running ZITADEL</li>
            <li>Running backups for ZITADEL</li>
        </ul>
        <p>Time Estimate</p>
        <ul>
            <li>8 hours</li>
        </ul>
    </TabItem>
    <TabItem value="orbos">
        <p>In Scope</p>
        <ul>
            <li>Check prerequisites and architecture</li>
            <li>Setup of Kubernetes cluster with ORBOS</li>
            <li>Setup of in-cluster toolset with ORBOS, which includes monitoring and log collecting</li>
            <li>Installation and configuration of ZITADEL with the ZITADEL-operator</li>
            <li>Installation and configuration of CockroachDB with the Database-operator</li>
            <li>Functional testing of the ZITADEL instance</li>
        </ul>
        <p>Out of Scope</p>
        <ul>
            <li>Integration of external S3-storage or other types of storage</li>
            <li>Integration into internal monitoring and alerting</li>
            <li>Multi-cluster architecture deployments</li>
            <li>Changes for specific environments</li>
            <li>Performance testing</li>
            <li>Production deployment</li>
            <li>Application-side coding, configuration, or tuning</li>
            <li>Changes or configuration on assets used in ZITADEL</li>
            <li>Setting up or maintaining backup</li>
        </ul>
        <p>Prerequisites</p>
        <ul>
            <li>Possibility to run ORBOS with one of the implemented Providers</li>
            <li>S3-storage for assets in ZITADEL</li>
            <li>S3-storage for backups of the database</li>
            <li>CDN or internal DNS to create DNS records for 4 subdomains pointing at ZITADEL</li>
            <li>Email relay to send emails from ZITADEL</li>
            <li>Account for Twilio to send SMS from ZITADEL</li>
        </ul>
        <p>Deliverable</p>
        <ul>
            <li>Running Vanilla Kubernetes</li>
            <li>Running toolset for monitoring and alerting</li>
            <li>Running CockroachDB</li>
            <li>Running ZITADEL</li>
            <li>Running backups for ZITADEL</li>
        </ul>
        <p>Time Estimate</p>
        <ul>
            <li>12 hours</li>
        </ul>
    </TabItem>
    <TabItem value="checkup">
        <p>In Scope</p>
        <ul>
            <li>Check prerequisites and architecture</li>
            <li>Check configuration for ZITADEL and ORBOS</li>
            <li>Functional testing of the ZITADEL instance</li>
        </ul>
        <p>Out of Scope</p>
        <ul>
            <li>Integration of external S3-storage or other types of storage</li>
            <li>Integration into internal monitoring and alerting</li>
            <li>Changes for specific environments</li>
            <li>Performance testing</li>
            <li>Application-side coding, configuration, or tuning</li>
            <li>Changes or configuration on assets used in ZITADEL</li>
            <li>Setting up or maintaining backup</li>
        </ul>
        <p>Prerequisites</p>
        <ul>
            <li>Access to the Kubernetes cluster or physical/virtual nodes</li>
            <li>Environment to test and check should not be in productive use</li>
            <li>Access to used configuration</li>
        </ul>
        <p>Deliverable</p>
        <ul>
            <li>Document detailing findings and description of the suggested configuration changes</li>
        </ul>
        <p>Time Estimate</p>
        <ul>
            <li>10 hours</li>
        </ul>
    </TabItem>
</Tabs>