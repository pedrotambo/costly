// import { group } from "console";
// import { useListGroupsQuery } from "../../services/api";
// import { Link } from "react-router-dom";
// import { Card, CardHeader, Heading } from "@chakra-ui/react";

export const Home = () => {
//   const { data, error, isLoading } = useListGroupsQuery();
  console.log("running...")
  const data = 1

  return (
    <div className="w-1/2">
      <p>No hay grupos al momento...</p>
      {!data && (
        <p>No hay grupos al momento...</p>
      )}
      <ul>
        {
        //   !isLoading && data &&
        //   data.map(group => (
        //     <li>
        //       <Link
        //         to={`/groups/${group.id}`}
        //       >
        //         <Card
        //           align='center'
        //           variant='filled'
        //         >
        //           <CardHeader>
        //             <Heading size='md'>
        //               {group.name}
        //             </Heading>
        //           </CardHeader>
        //         </Card>
        //       </Link>
        //     </li>
        //   ))
        }
      </ul>
    </div>
  );
}
