import {useEffect, useState} from "react";
import {ResourceProps} from "@refinedev/core";
import {RouteProps} from "react-router-dom";
import {gql, GraphQLClient} from "@refinedev/graphql";
import {GetDefaultActionConfig} from "./xgen_graphql_defaults";

type XgenSet = {
    resources: ResourceProps[],
    routes: RouteProps[],
}

export function useXgenSet(client: GraphQLClient): XgenSet {
    const [resources, setResources] = useState<ResourceProps[]>([]);
    const [routes, setRoutes] = useState<RouteProps[]>([]);

    useEffect(() => {
        client.request(gql`
            query{
                _xgen_introspection{
                    _per_def{
                        XgenResource{
                            name
                            value{
                                Route
                                Name
                            }
                        }
                        XgenResourceAction{
                            name,
                            value{
                                Resource
                                Action
                                Route
                            }
                        }
                        XgenResourceListAction{
                            name,
                            value{
                                Resource
                                Action
                                Route
                            }
                        }
                    }
                }
            }
        `).then(async (data: any) => {
            const _resources: Record<string, ResourceProps> = {};
            const _routes: RouteProps[] = []

            for (const action of data._xgen_introspection._per_def.XgenResourceListAction) {

                const resource = data._xgen_introspection._per_def.XgenResource.find((r: any) => r.value.Name === action.value.Resource)
                if (!resource) {
                    console.log(`No resource for ${action.value.Resource}`)
                    continue;
                }
                const resourceRoute = resource.value.Route || resource.value.Name

                const defaultConfig = await GetDefaultActionConfig(client, action.name, action.value.Resource, action.value.Action)
                if (!defaultConfig) {
                    console.log(`No default config for ${action.value.Resource} ${action.value.Action}`)
                    continue;
                }
                const actionRoute = action.value.Route || defaultConfig.pattern

                const finalRoute = [resourceRoute, actionRoute].join("/").replace(/\/+$/g, "")

                _routes.push({
                    index: Boolean(resource.value.Primary),
                    path: finalRoute,
                    element: defaultConfig.element,
                })
                _resources[action.value.Resource] = {
                    ..._resources[action.value.Resource],
                    name: action.value.Resource,
                    [defaultConfig.action]: finalRoute,
                }
            }
            setResources(() => Object.values(_resources));
            setRoutes(() => _routes);
        })
    }, [client])

    return {
        resources,
        routes,
    }
}