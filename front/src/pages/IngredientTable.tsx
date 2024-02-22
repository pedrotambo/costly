import { Ingredient, useGetIngredientsQuery } from "../services/api";
import { createColumnHelper } from "@tanstack/react-table"
import { DataTable } from "../components/Table/DataTable"

export const IngredientTable = () => {
    const { data, error, isLoading } = useGetIngredientsQuery();
    const columnHelper = createColumnHelper<Ingredient>();

    const columns = [
    columnHelper.accessor("id", {
			cell: (info) => info.getValue(),
			header: "ID"
    }),
    columnHelper.accessor("name", {
			cell: (info) => info.getValue(),
			header: "Name"
    }),
    columnHelper.accessor("price", {
			cell: (info) => info.getValue(),
			header: "Price",
			meta: {
				isNumeric: true
			},
    }),
    columnHelper.accessor("unit", {
        cell: (info) => info.getValue(),
        header: "Unit",
        meta: {
					isNumeric: true
        }
    })
    ];
    return (
        <DataTable columns={columns} data={data ?? []} />
    )
}
