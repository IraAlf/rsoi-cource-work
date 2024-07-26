import { Box } from "@chakra-ui/react";
import RecipeMap from "components/RecipeMap";
import React from "react";
import { NavigateFunction, Params } from "react-router-dom";

import styles from "./CategoryPage.module.scss";

interface CategoryProps {
    match: Readonly<Params<string>>
    navigate: NavigateFunction
    cookie: {
        role?: string;
        login?: string;
    }
}

const CategoryPage: React.FC<CategoryProps> = (props) => {

  return (
    <Box className={styles.category_page}>
      <RecipeMap searchQuery={props.match.title}/>
    </Box>
  );
};

export default CategoryPage;
