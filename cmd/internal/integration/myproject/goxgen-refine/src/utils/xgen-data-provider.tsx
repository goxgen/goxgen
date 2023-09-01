import {
    BaseRecord,
    CrudFilters,
    CrudSorting,
    DataProvider,
    LogicalFilter,
} from "@refinedev/core";
import { GraphQLClient } from "graphql-request";
import * as gql from "gql-query-builder";
// import pluralize from "pluralize";
import camelCase from "camelcase";
// import * as SorterEnums from "./sorter-enums"

export type Cursors = {
    before?: string
    after?: string
    first?: number
    last?: number
}

export type OrderQuery = {
    field: string,
    direction: "ASC" | "DESC"
}

function camel(str: string): string{
    return camelCase(str, {preserveConsecutiveUppercase: true})
}
function pascal(str: string): string{
    return camelCase(str, {preserveConsecutiveUppercase: true, pascalCase: true})
}

export const generateOrderQuery = (resource: string, sort?: CrudSorting): undefined|OrderQuery => {
    // if (sort && sort.length > 0) {
    //
    //     if (sort.length > 1) {
    //         console.log("Multiple sorts not supporting because of EntGQL right not supporting multiple orders, only first order is working: https://github.com/ent/ent/issues/2198")
    //     }
    //
    //     const fieldGroupName :keyof typeof SorterEnums = pascal(resource)+"SorterEnums" as any
    //     const fieldGroup = SorterEnums[fieldGroupName] as any
    //     return {
    //         field: fieldGroup[sort[0].field],
    //         direction: sort[0].order.toUpperCase() as "ASC"|"DESC",
    //     }
    // }
    return undefined;
};

export const generateFilter = (filters?: CrudFilters): Record<string, any> => {
    const queryFilters: Record<string, any>  = {};

    if (filters) {
        filters.map((filter) => {
            if (
                filter.operator !== "or" &&
                filter.operator !== "and" &&
                "field" in filter
            ) {
                const { field, operator, value } = filter;
                if (operator === "eq") {
                    queryFilters[camel(field)] = value;
                } else {
                    queryFilters[camel(`${field}_${operator}`)] = value;
                }
            } else {
                const value = filter.value as LogicalFilter[];
                const orFilters: any[] = [];
                value.map((val) => {
                    if (val.operator === "eq") {
                        orFilters.push({
                            [ camel(`${val.field}`)]: val.value,
                        });
                    } else {
                        orFilters.push({
                            [ camel(`${val.field}_${val.operator}`)]: val.value,
                        });
                    }
                });
                queryFilters["or"] = orFilters;
            }
        });
    }

    return queryFilters;
};

export default (client: GraphQLClient): Required<DataProvider> => {
    return {
        /** Done */
        getOne: async ({ resource, id, meta }) => {

            return {
                data: {
                    id: id,
                    name: "test",
                },
            } as any
        },
        /** Done */
        getList: async ({
                            resource,
                            pagination = true,
                            sorters,
                            filters,
                            meta,
                        }) => {
            console.log(resource)
            const {
                current = 1,
                pageSize = 10,
                mode = "server",
            } = pagination ?? {};

            // const sortBy = generateSort(sorters);
            const filterBy = generateFilter(filters);

            const camelResource = camelCase(resource);

            const operation = resource+"_browse";

            const { query, variables } = gql.query({
                operation,
                variables: {
                    // ...meta?.variables,
                    // sort: sortBy,
                    // where: { value: filterBy, type: "JSON" },
                    // ...(mode === "server"
                    //     ? {
                    //         start: (current - 1) * pageSize,
                    //         limit: pageSize,
                    //     }
                    //     : {}),
                },
                fields: meta?.fields,
            });

            const response = await client.request<BaseRecord>(query, variables);

            return {
                data: response[operation],
                total: response[operation].count,
            };


            return {
                data: [
                    {
                        id: "1",
                    }
                ],
                pageInfo: {},
                total: 10,
            };
        },

        getMany: async ({ resource, ids, meta }) => {
            return {
                data: [],
            };
        },

        create: async ({ resource, variables, meta }) => {
            return {
                data: {
                    id: "1",
                },
            } as any;
        },

        createMany: async ({ resource, variables, metaData }) => {
            return {
                data: [],
            }
        },

        update: async({ resource, id, variables, meta }) => {
            return {
                data: {
                    id: id,
                },
            } as any;
        },

        updateMany: async ({ resource, ids, variables, meta }) => {
            return {
                data: [],
            };
        },



        deleteOne: async ({ resource, id, meta }) => {
            return {
                data: {
                    id : id,
                },
            } as any;
        },

        deleteMany: async ({ resource, ids, metaData }) => {
            return {
                data: [],
            };
        },

        getApiUrl: () => {
            throw Error("Not implemented on refine-graphql data provider.");
        },

        custom: async ({ url, method, headers, metaData }) => {
            return {data: {}} as any
        }
    };
};