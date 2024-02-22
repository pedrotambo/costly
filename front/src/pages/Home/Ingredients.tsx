import { useGetIngredientsQuery } from "../../services/api";
import { Link } from "react-router-dom";
// import { Card, CardHeader, Heading } from "@chakra-ui/react";
import {
  Table,
  Thead,
  Tbody,
  Tfoot,
  Tr,
  Th,
  Td,
  TableCaption,
  TableContainer,
  Tooltip,
  Button,
  HStack,
  Heading,
  ButtonGroup,
  // Input
} from '@chakra-ui/react'

import { AddIcon, SearchIcon, ChevronDownIcon, ChevronUpIcon } from '@chakra-ui/icons'

function firstLetterUpperCase(word: string) {
  return word[0].toUpperCase() + word.substring(1, word.length)
}

// function TableHeader({ onAdd: React.MouseEventHandler<HTMLButtonElement>, onSearch, onSort }) {
function TableHeader() {
  return (
    <div>
      <Heading as="h2" size="md" marginBottom="4">Table Title</Heading>
      <ButtonGroup marginBottom="4">
        <Button leftIcon={<AddIcon />}>Add New Item</Button>
        {/* <Input placeholder="Search" onChange={onSearch} marginRight="4" /> */}
        <Button leftIcon={<SearchIcon />}>Search</Button>
        <Button leftIcon={<ChevronDownIcon />}>Sort Asc</Button>
        <Button leftIcon={<ChevronUpIcon />}>Sort Desc</Button>
      </ButtonGroup>
    </div>
  );
}

  // const handleAddNewItem = () => {
  //   // Add logic to add new item
  // };

  // const handleSearch = (event) => {
  //   const searchTerm = event.target.value;
  //   // Add logic to search items
  // };

  // const handleSort = (order) => {
  //   // Add logic to sort items
  // };


export const Ingredients = () => {
  const { data, error, isLoading } = useGetIngredientsQuery();

  // return (
  //   <div className="w-1/2">
  //     {!data && (
  //       <p>No hay grupos al momento...</p>
  //     )}
  //     <ul>
  //       {
  //         !isLoading && data &&
  //         data.map(ingredient => (
  //           <li>
  //             <Link
  //               to={`/groups/${ingredient.id}`}
  //             >
  //               <Card
  //                 align='center'
  //                 variant='filled'
  //               >
  //                 <CardHeader>
  //                   <Heading size='md'>
  //                     {ingredient.name} ${ingredient.price}/{ingredient.unit}
  //                   </Heading>
  //                 </CardHeader>
  //               </Card>
  //             </Link>
  //           </li>
  //         ))
  //       }
  //     </ul>
  //   </div>
  // );
  const header = TableHeader()
  return (
    <div>
    {/* <Tooltip label='Hover me' placement='right-end'>
      Ingredients
      <Button>asdf</Button>
    </Tooltip> */}
    {header}
    <HStack spacing={6}>
      INGREDIENTS
      {/* <Tooltip label='Top start' placement='top-start'>
        <Button>Top-Start</Button>
      </Tooltip>

      <Tooltip label='Top' placement='top'>
        <Button>Top</Button>
      </Tooltip> */}

      <Tooltip label='Right end' placement='right-end'>
        <Button>Right-End</Button>
      </Tooltip>
    </HStack>
    <TableContainer>
      <Table variant='simple'>
        <TableCaption>List of ingredients with its prices per unit</TableCaption>
        <Thead>
          <Tr>
            <Th>ID</Th>
            <Th>Name</Th>
            <Th>Price</Th>
            {/* <Th isNumeric>Unit</Th> */}
          </Tr>
        </Thead>
        <Tbody>
          {
            data?.map(ingredient => {
              return (
                <Tr>
                  <Td>{ingredient.id}</Td>
                  <Td>{firstLetterUpperCase(ingredient.name)}</Td>
                  <Td>${ingredient.price}/{ingredient.unit}</Td>
                </Tr>
              )
            })
          }
        </Tbody>
        {/* <Tfoot>
          <Tr>
            <Th>Ingredients</Th>
            <Th>Price</Th>
            <Th isNumeric>Unit</Th>
          </Tr>
        </Tfoot> */}
      </Table>
    </TableContainer>
    </div>
  )
}
    

  