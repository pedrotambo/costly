import { Recipe, useGetRecipesQuery } from "../services/api";
import { createColumnHelper } from "@tanstack/react-table"
import { DataTable } from "../components/Table/DataTable"

export const RecipeTable = () => {
    const { data, error, isLoading } = useGetRecipesQuery();
    const columnHelper = createColumnHelper<Recipe>();

    const columns = [
    columnHelper.accessor("id", {
			cell: (info) => info.getValue(),
			header: "ID"
    }),
    columnHelper.accessor("name", {
			cell: (info) => info.getValue(),
			header: "Name"
    }),
    columnHelper.accessor("ingredients", {
				cell: ({ row }) => `${row.original.ingredients.map(i => `${i.units}${i.ingredient.unit} of ${i.ingredient.name}`).join(', ')}`,
        header: "Ingredients",
				enableSorting: false
    }),
		columnHelper.accessor("cost", {
			cell: (info) => info.getValue(),
			header: "Cost",
    }),
    ];
    return (
        <DataTable columns={columns} data={data ?? []} />
    )
}
