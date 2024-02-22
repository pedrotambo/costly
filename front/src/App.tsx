import './App.css';
import { Outlet } from 'react-router-dom';
import { Avatar, HStack } from '@chakra-ui/react';

function App() {
  return (
    <div className='w-full justify-center flex flex-col items-center'>
      <Nav></Nav>
      <Outlet></Outlet>
    </div>
  );
}

const Nav = () => {
  return (
    <HStack
      className='w-full px-3 py-1 border-b border-gray-200 mb-2'
      justifyContent="flex-start"
      alignItems="center" // Add this line to align vertically to the middle
    >
      <h1 className='w-full text-2xl font-bold'>
        Costly
      </h1>
      <Avatar
        name="Pedro Tamborindeguy"
        size="sm"
      >
      </Avatar>
    </HStack>
  )
};

export default App;
