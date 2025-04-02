import { useQuery } from "@tanstack/react-query";
import { BASE_URL } from "../config";
import { Box, Card, CardContent, Typography } from "@mui/material";

export type Habit = {
    _id: number;
    user_id: number;
    name: string;
    description?: string;
    frequency: "daily" | "weekly" | "monthly";
    remind_time: string;
    timezone: string;
    created_at: string;
    updated_at: string;
}

function HomePage() {
    const { data: habits, isLoading } = useQuery<Habit[]>({
        queryKey: ["habits"],
        queryFn: async () => {
            try {
                const res = await fetch(BASE_URL + "/habits")

                const contentType = res.headers.get("content-type");
                if (!contentType || !contentType.includes("application/json")) {
                    throw new Error("Server did not return JSON");
                }

                const data = await res.json();

                if (!res.ok) {
                    throw new Error(data.error || "Something went wrong")
                }

                return data || []
            } catch (error) {
                console.log(error)
            }
        }
    });

    if (isLoading) { return (<Typography variant="h2" sx={{ textAlign: 'center' }}>Loading...</Typography>) }

    return (
        <Box>
            {habits?.map((habit: Habit) => (
                <Card variant="outlined" key={habit._id}>
                    <CardContent>
                        <Typography sx={{ color: 'text.secondary', fontSize: 14 }}>
                            User id: <b>{habit.user_id}</b>
                        </Typography>
                    </CardContent>
                    <CardContent>
                        <Typography variant="h3" component="div">
                            {habit.name}
                        </Typography>
                    </CardContent>
                    <CardContent>
                        <Typography variant="body1">
                            {habit.description || "There is no description here("}
                        </Typography>
                    </CardContent>
                    <CardContent>
                        <Typography variant="h6">
                            Rimend at: {habit.remind_time} ({habit.frequency})
                        </Typography>
                    </CardContent>
                    <CardContent>
                        <Typography variant="body2" sx={{ color: 'text.secondary', fontSize: 12 }}>Timezone: <b>{habit.timezone}</b></Typography>
                    </CardContent>
                </Card>
            ))}
        </Box>
    )
}

export default HomePage;