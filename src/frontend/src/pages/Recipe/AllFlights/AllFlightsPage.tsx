import { Box } from "@chakra-ui/react";
import CategoryMap from "components/CategoryMap";
import { SearchContext } from "context/Search";
import GetCategories from "postAPI/categories/GetAll";
import React, { useContext, useEffect } from "react";
import RecipeMap from "../../../components/RecipeMap/RecipeMap";

import styles from "./AllFlightsPage.module.scss";
import GetFlights from "postAPI/flights/GetAll";

interface AllRecipesProps {}

const AllFlightsPage: React.FC<AllRecipesProps> = (props) => {
  const searchContext = useContext(SearchContext);

  useEffect(() => {
    console.log("AllFlightsPage mounted");
    console.log("Search query:", searchContext.query);
  }, [searchContext]);

  return (
    <Box className={styles.main_box}>
      <RecipeMap searchQuery={searchContext.query} getCall={GetFlights}/>
    </Box>
  );
};

export default AllFlightsPage;
