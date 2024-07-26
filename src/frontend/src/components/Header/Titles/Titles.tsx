import React from "react";
import {
  Text,
  Box,
} from "@chakra-ui/react";

import styles from "./Titles.module.scss";


export interface TitlesProps {
  title: string
  undertitle?: string
}

const Titles: React.FC<TitlesProps> = (props) => {
  return (

    <Box className={styles.titles}>
        <Text id="title" className={styles.title}>
            {props.title}
        </Text>
        <Text id="undertitle" className={styles.undertitle}>
            {props.undertitle ? props.undertitle : "ㅤ"}
        </Text>
    </Box>
  );
};

export default React.memo(Titles);
