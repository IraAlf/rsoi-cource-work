import { Box } from "@chakra-ui/react";
import React from "react";
import { useCookies } from "react-cookie";
import { useParams } from "react-router-dom";
import RecipeMap from "../../../components/RecipeMap/RecipeMap";

import styles from "./LikedRecipesPage.module.scss";

interface LikedRecipesProps {}

const LikedRecipes: React.FC<LikedRecipesProps> = (props) => {
  let [cookie] = useCookies(["role", "login"]);
  const params = useParams();
  let login = params.login ? params.login : cookie.login;

  return (
    <Box className={styles.main_box}>
      <RecipeMap getCall={() => {}} />
    </Box>
  );
};

export default LikedRecipes;
