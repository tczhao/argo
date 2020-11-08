import {SlideContents} from 'argo-ui';
import * as React from 'react';

import * as models from '../../../../models';
import {ResourceEditor} from '../../../shared/components/resource-editor/resource-editor';
import {getResolvedTemplates} from '../../../shared/template-resolution';

export interface WorkflowYamlViewerProps {
    workflow: models.Workflow;
    selectedNode: models.NodeStatus;
}

export class WorkflowYamlViewer extends React.Component<WorkflowYamlViewerProps> {

    public render() {
        const contents: JSX.Element[] = [];
        contents.push(<h3 key='title'>Node</h3>);
        if (this.props.selectedNode) {
            const parentNode = this.props.workflow.status.nodes[this.props.selectedNode.boundaryID];
            if (parentNode) {
                contents.push(
                    <div key='parent-node'>
                        <h4>{this.normalizeNodeName(this.props.selectedNode.displayName || this.props.selectedNode.name)}</h4>
                        <ResourceEditor kind='Template' value={getResolvedTemplates(this.props.workflow, parentNode)} />
                    </div>
                );
            }

            contents.push(
                <div key='current-node'>
                    <h4>{this.props.selectedNode.name}</h4>
                    <ResourceEditor kind='Template' value={getResolvedTemplates(this.props.workflow, this.props.selectedNode)} />
                </div>
            );
        }
        const templates = this.props.workflow.spec.templates;
        if (templates && Object.keys(templates).length) {
            contents.push(<SlideContents title='Templates' key='templates' contents={<ResourceEditor kind='Template' value={templates} />} className='workflow-yaml-section' />);
        }
        const storedTemplates = this.props.workflow.status.storedTemplates;
        if (storedTemplates && Object.keys(storedTemplates).length) {
            contents.push(
                <SlideContents
                    title='Stored Templates'
                    key='stored-templates'
                    contents={<ResourceEditor kind='Template' value={storedTemplates} />}
                    className='workflow-yaml-section'
                />
            );
        }

        return (
            <div className='workflow-yaml-viewer'>
                {contents}
            </div>
        );
    }

    private normalizeNodeName(name: string) {
        const parts = name.replace(/([(][^)]*[)])/g, '').split('.');
        return parts[parts.length - 1];
    }

}
