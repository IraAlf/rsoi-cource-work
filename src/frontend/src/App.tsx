import * as React from "react";
import theme from "./styles/extendTheme";
import { ChakraProvider, Box, Container } from "@chakra-ui/react";
import { Routes, Route } from "react-router";
import { BrowserRouter } from "react-router-dom";

import Login from "pages/Login";
import SignUp from "pages/Signup";
import RecipeInfoPage from "pages/Recipe/RecipeInfo";
import UsersPage from "pages/Users/UsersPage";
import AuthorRecipesPage from "pages/Recipe/AuthorRecipes";
import LikedRecipesPage from "pages/Recipe/LikedRecipes";
import CategoryPage from "pages/Category";

import SearchContextProvider from "context/Search";
import { HeaderRouter } from "components/Header";
import AllFlightsPage from "pages/Recipe/AllFlights/AllFlightsPage";
import FlightInfoPage from "pages/Recipe/FlightInfo";
import AllTicketsPage from "pages/Ticket/AllTickets";
import StatisticsPage from "pages/Statistics/StatisticsPage";


interface HomeProps { }
const Home: React.FC<HomeProps> = () => {
  return (
    <Box backgroundColor="bg" h="auto">
      <Container maxW="1000px" minH="95%"
        display="flex"
        paddingX="0px" paddingY="30px"
        alignSelf="center" justifyContent="center"
        textStyle="body"
      >
        <Routing />
      </Container>
    </Box>
  );
};

function Routing() {
  return <BrowserRouter>
    <Routes>
      <Route path="/" element={<Login />} />
      <Route path="/auth/signin" element={<Login />} />
      <Route path="/auth/signup" element={<SignUp />} />
      <Route path="/flights/:flightNumber" element={<FlightInfoPage />} />
      <Route path="/tickets" element={<AllTicketsPage />} />
      <Route path="/statistics" element={<StatisticsPage />} />
      <Route path="/flights" element={<AllFlightsPage />} />



      <Route path="*" element={<NotFound />} />
    </Routes>
  </BrowserRouter>
}

function NotFound() {
  return <h1>Page not Found</h1>
}

export const App = () => {
  return (
    <ChakraProvider theme={theme}>
      <SearchContextProvider>
        <HeaderRouter />
        <Home />
      </SearchContextProvider>
    </ChakraProvider>
  )
};
