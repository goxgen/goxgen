import {DataGrid, GridColDef} from "@mui/x-data-grid";
import {IResourceComponentsProps, useDelete, useNavigation, useTranslate,} from "@refinedev/core";
import {List, useDataGrid} from "@refinedev/mui";
import React, {useEffect} from "react";

// import { SomeInterface } from "interfaces";
import {GraphQLClient} from "@refinedev/graphql";

import * as gql from "gql-query-builder";

interface SomeInterface {
    name: string;
    surname: string;
    avatar: string;
    address: string;
    id: number;
}

type XgenListComponentProps = {
    client: GraphQLClient,
    objectName: string,
    action: string
}
export const XgenListComponent:
React.FC<IResourceComponentsProps & XgenListComponentProps> = ({client, objectName}) => {
    const [columns, setColumns] = React.useState<GridColDef[]>([]);
    const [fields, setFields] = React.useState<string[]>([]);
    useEffect(() => {
        const {query, variables} = gql.query({
            operation: "_xgen_introspection",
            fields: [
                {
                    _per_object: [
                        {
                            [objectName]: [
                                {
                                    object:[
                                        {
                                            XgenResource: [
                                                "Name"
                                            ]
                                        }
                                    ]
                                },
                                {
                                    field: [
                                        "name",
                                        {
                                            definition: [
                                                {
                                                    XgenResourceActionField: [
                                                        "Label",
                                                        "Description"
                                                    ]
                                                }
                                            ]
                                        }
                                    ]
                                }
                            ]
                        }
                    ]
                }
            ]
        })
        console.log("asdasd",query, variables)
        // client.request<any>(gql`
        //             query{
        //                 _xgen_introspection {
        //                     _per_object {
        //                         ${objectName} {
        //                 object {
        //                     XgenResource {
        //                         Name
        //                     }
        //                 }
        //                 field {
        //                     name
        //                     definition {
        //                         XgenResourceActionField {
        //                             Label
        //                             Description
        //                         }
        //                     }
        //                 }
        //             }
        //             }
        //             }
        //             }
        //         `)
        client.request<any>(query, variables).then((data: any) => {
            const fieldsDefs = data._xgen_introspection._per_object[objectName].field
            const _fields = fieldsDefs.map((field: any) => {
                return {
                    field: field.name,
                    headerName: field.definition?.XgenResourceActionField.Label || field.name,
                    flex: 1,
                    minWidth: 200,
                }
            })
            setColumns(() => _fields)
            setFields(() => fieldsDefs.map((field: any) => field.name))
        })
    }, [client, objectName]);

    const {show, edit} = useNavigation();
    const {mutate: mutateDelete} = useDelete();

    const {dataGridProps} = useDataGrid<SomeInterface>({
        initialPageSize: 10,
        initialSorter: [
            {
                field: "id",
                order: "desc",
            },
        ],
        meta: {
            fields: fields,
        }
    });
    const t = useTranslate();


    return <List wrapperProps={{sx: {paddingX: {xs: 2, md: 0}}}}>
        <DataGrid
            {...dataGridProps}
            columns={columns}
            autoHeight
            pageSizeOptions={[10, 20, 50, 100]}
            density="comfortable"
            sx={{
                "& .MuiDataGrid-cell:hover": {
                    cursor: "pointer",
                },
            }}
            onRowClick={(row) => {
                show("couriers", row.id);
            }}
        />
    </List>

}