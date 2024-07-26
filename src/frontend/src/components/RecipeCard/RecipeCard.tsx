import { Box, HStack, Link, Text, VStack } from "@chakra-ui/react";
import React from "react";

import { Flight as FlightI } from "types/Flight";

import styles from "./RecipeCard.module.scss";

interface FlightProps extends FlightI {}

const RecipeCard: React.FC<FlightProps> = (props) => {
  var path = "/flights/" + props.flightNumber;

  return (
    <Link className={styles.link_div} href={path}>
      <Box className={styles.main_box}>
        <Box className={styles.info_box}>
          <VStack>
            <Box className={styles.description_box}>
              <Text>{props.flightNumber}</Text>
            </Box>
            <Box className={styles.title_box}>
              <Text>{props.fromAirport}➤</Text>
            </Box>
            <Box className={styles.title_box}>
              <Text>➤{props.toAirport}</Text>
            </Box>

            <HStack>
              <Box className={styles.description_box}>
                <Text>{props.date}</Text>
              </Box>
              <Box className={styles.description_box}>
                <Text>{props.price}</Text>
              </Box>
            </HStack>
          </VStack>
        </Box>
      </Box>
    </Link>
  );
};

export default RecipeCard;
