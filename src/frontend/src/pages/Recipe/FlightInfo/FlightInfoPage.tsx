import { VStack, Text, Box, Switch } from "@chakra-ui/react"
import GetFlight from "postAPI/flights/Get"
import React from "react"
import { NavigateFunction, Params } from "react-router-dom"
import { Flight as FlightT} from "types/Flight"
import styles from "./FlightsInfoPage.module.scss";
import RoundButton from "components/RoundButton/RoundButton"

import CreateTicket from "postAPI/tickets/Create"

type State = {
    flight?: FlightT
    useBonusPoints: boolean
}

type RecipeInfoParams = {
    match: Readonly<Params<string>>
    navigate: NavigateFunction
}

class FlightInfoPage extends React.Component<RecipeInfoParams, State> {
    flightNumber: string;

    constructor(props) {
        super(props);
        this.flightNumber = this.props.match.flightNumber || "?";
        this.state = {
            useBonusPoints: false
        }
    };

    componentDidMount(): void {
        GetFlight(this.flightNumber).then(data => {
            console.log(data)
            if (data.status === 200) {
                this.setState({flight: data.content})
            }
        });
    }


    airport_img() {
        return (
            <img
                src="/airport.svg"
                className={styles.airports_img}
            />
        )
    }

    fromBox() {
        return (<Box >
            {this.airport_img()}
            <Text>{this.state.flight?.fromAirport}</Text>
        </Box>);
    };

    toBox() {
        return (<Box>
            {this.airport_img()}
            <Text>{this.state.flight?.toAirport}</Text>
        </Box>);
    };

    submit(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
        if (!this.state.flight) {
            console.error("flight is null")
            return
        }

        let button = e.currentTarget
        button.disabled = true
        CreateTicket(this.state.flight.flightNumber, this.state.flight.price, this.state.useBonusPoints)
        .then(data => {
            button.disabled = false
            if (data.status === 200) {
                window.location.href = '/tickets';
            } else {
                var title = document.getElementById("undertitle")
                if (title)
                    title.innerText = "Ошибка покупки билета"
            }
        });
    }

    async handleToggle() {
        this.setState({useBonusPoints: !this.state.useBonusPoints})
    }

    render() {
        return (
            <VStack className={styles.main_box}>
                {this.state.flight &&
                    <Box style={{width: "100%"}}>
                        <Text>Рейс номер {this.flightNumber}</Text>
                        <Box className={styles.airports_box}>
                            {this.fromBox()}
                            <Text>{this.state.flight?.date}</Text>
                            {this.toBox()}
                        </Box>

                        <Box className={styles.bonus_div}>
                            <Switch
                                isChecked={this.state.useBonusPoints}
                                onChange={() => this.handleToggle()}
                                colorScheme="teal"
                                size="md"
                            />
                            <Text>
                                Использовать бонусные баллы
                            </Text>
                        </Box>


                        <RoundButton className={styles.buy_button} type="submit" onClick={event => this.submit(event)}>
                            Забронировать билет за {this.state.flight?.price}
                        </RoundButton>
                    </Box>
                }
            </VStack>
        );
    };
};

export default FlightInfoPage;