import { Box } from "@chakra-ui/react";
import React from "react";
import RecipeCard from "../RecipeCard";
import { AllFilghtsResp } from "postAPI"

import styles from "./RecipeMap.module.scss";

interface RecipeBoxProps {
    searchQuery?: string
    getCall: (page: number, size: number) => Promise<AllFilghtsResp>
}

type State = AllFilghtsResp

class RecipeMap extends React.Component<RecipeBoxProps, State> {
    constructor(props) {
        super(props);
    }

    async getAll() {
        var data = await this.props.getCall(0, 20)
        if (data)
            this.setState(data)
    }

    componentDidMount() {
        this.getAll()
    }

    componentDidUpdate(prevProps) {
        if (this.props.searchQuery !== prevProps.searchQuery) {
            this.getAll()
        }
    }

    render() {
        return (
            <Box className={styles.map_box}>
                {this.state?.items.map(item => <RecipeCard {...item} key={item.flightNumber}/>)}
            </Box>
        )
    }
}

export default React.memo(RecipeMap);
