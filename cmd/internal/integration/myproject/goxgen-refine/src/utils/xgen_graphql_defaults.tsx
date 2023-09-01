import {ReactElement} from "react";
import {List} from "@refinedev/mui";
import {GraphQLClient} from "@refinedev/graphql";
import {XgenListComponent} from "./xgen-list";

export type DefaultActionConfig = {
    action: string,
    pattern: string,
    element: ReactElement,
}

export async function GetDefaultActionConfig(
    client: GraphQLClient,
    objectName: string,
    resource: string,
    action: string,
): Promise<DefaultActionConfig|undefined> {
    const defs: Record<string, DefaultActionConfig> =  {
        BROWSE_QUERY: {
            action : 'list',
            pattern: ``,
            element: <XgenListComponent client={client} objectName={objectName} action={action}></XgenListComponent>
        },
        SHOW_QUERY: {
            action : 'show',
            pattern: `:id`, element: <div>{action}</div>},
        UPDATE_MUTATION: {
            action : 'edit',
            pattern:`:id/edit`, element: <div>{action}</div>},
        CREATE_MUTATION: {
            action : 'create',
            pattern:`create`, element: <div>{action}</div>},
        DELETE_MUTATION: {
            action : 'delete',
            pattern:`:id/delete`, element: <div>{action}</div>},
    }
    return defs[action]
}