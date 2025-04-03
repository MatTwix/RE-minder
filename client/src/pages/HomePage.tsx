import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "../config";
import { Box, Button, Card, CardContent, IconButton, MenuItem, TextField, Typography } from "@mui/material";
import { useState } from "react";
import { Delete, Update } from "@mui/icons-material";

export type Habit = {
    id?: number;
    user_id: number;
    name: string;
    description?: string;
    frequency: "daily" | "weekly" | "monthly";
    remind_time: string;
    timezone: string;
    created_at?: string;
    updated_at?: string;
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
            } catch (error: any) {
                throw new Error(error)
            }
        }
    });

    const queryClient = useQueryClient();
    const [createUserId, setCreateUserId] = useState('')
    const [createName, setCreateName] = useState('')
    const [createDescription, setCreateDescription] = useState('')
    const [createRemindTime, setCreateRemindTime] = useState('')
    const [createFrequency, setCreateFrequency] = useState('')
    const [createTimezone, setCreateTimezone] = useState('')

    const { mutate: createHabit, isPending: isCreating } = useMutation({
        mutationKey: ["createHabit"],
        mutationFn: async (e: React.FormEvent) => {
            e.preventDefault();
            if (!createUserId || !createName || !createRemindTime || !createFrequency) return alert("Fill all required fields!")

            const habitData: Partial<Habit> = {
                user_id: Number(createUserId),
                name: createName,
                description: createDescription,
                remind_time: createRemindTime,
                frequency: createFrequency as "daily" | "weekly" | "monthly",
                timezone: createTimezone
            }

            try {
                const res = await fetch(BASE_URL + `/habits`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(habitData),
                });

                const contentType = res.headers.get("content-type");
                if (!contentType || !contentType.includes("application/json")) {
                    throw new Error("Server did not return JSON");
                }

                const data = await res.json();

                if (!res.ok) {
                    throw new Error(data.error || "Something went wrong!")
                }

                setCreateUserId("")
                setCreateName("")
                setCreateDescription("")
                setCreateRemindTime("")
                setCreateFrequency("")
                setCreateTimezone("")

                return data;
            } catch (error: any) {
                throw new Error(error);
            }
        }, 
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["habits"] })
        },
        onError: (error: any) => {
            alert(error.message)
        },
    })

    const [updateHabitId, setUpdateHabitId] = useState('')
    const [updateUserId, setUpdateUserId] = useState('')
    const [updateName, setUpdateName] = useState('')
    const [updateDescription, setUpdateDescription] = useState('')
    const [updateRemindTime, setUpdateRemindTime] = useState('')
    const [updateFrequency, setUpdateFrequency] = useState('')
    const [updateTimezone, setUpdateTimezone] = useState('')

    const { mutate: updateHabit, isPending: isUpdating } = useMutation({
        mutationKey: ["updateHabit"],
        mutationFn: async (e: React.FormEvent) => {
            e.preventDefault();
            if (!updateName || !updateRemindTime || !updateFrequency) return alert("Fill all required fields!")

            const habitData: Partial<Habit> = {
                user_id: Number(updateUserId),
                name: updateName,
                description: updateDescription,
                remind_time: updateRemindTime,
                frequency: updateFrequency as "daily" | "weekly" | "monthly",
                timezone: updateTimezone
            }

            try {
                const res = await fetch(`${BASE_URL}/habits/${updateHabitId}`, {
                    method: "PUT",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(habitData),
                });

                const contentType = res.headers.get("content-type");
                if (!contentType || !contentType.includes("application/json")) {
                    throw new Error("Server did not return JSON");
                }

                const data = await res.json();

                if (!res.ok) {
                    throw new Error(data.error || "Something went wrong!")
                }

                setUpdateHabitId("")
                setUpdateUserId("")
                setUpdateName("")
                setUpdateDescription("")
                setUpdateRemindTime("")
                setUpdateFrequency("")
                setUpdateTimezone("")

                return data;
            } catch (error: any) {
                throw new Error(error);
            }
        }, 
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["habits"] })
        },
        onError: (error: any) => {
            alert(error.message)
        },
    })

    const [deletingId, setDeletingId] = useState("")

    const { mutate:deleteHabit, isPending:isDeleting } = useMutation({
        mutationKey:["deleteHabit"],
        mutationFn: async () => {
            try {
                const res = await fetch(`${BASE_URL}/habits/${deletingId}`, {
                    method: "DELETE"
                })

                const contentType = res.headers.get("content-type");
                if (!contentType || !contentType.includes("application/json")) {
                    throw new Error("Server did not return JSON");
                }

                const data = await res.json();

                if (!res.ok) {
                    throw new Error(data.error || "Something went wrong")
                }

                console.log(data)

                return data || []
            } catch (error: any) {
                throw new Error(error)
            }
        },
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey:["habits"]})
        },
        onError: (error: any) => {
            alert(error.message)
        }
    })

    if (isLoading) { return (<Typography variant="h2" sx={{ textAlign: 'center' }}>Loading...</Typography>) }

    return (
        <Box sx={{ marginTop: 2 }}>
            <Box component="form" onSubmit={createHabit} sx={{ marginBottom: 2 }}>
                <Typography variant="h3">Create habit</Typography>
                <TextField 
                    required 
                    type="number" 
                    label="User_id" 
                    placeholder="Enter user_id" 
                    value={createUserId} 
                    onChange={(e) => setCreateUserId(e.target.value)} 
                />
                <TextField 
                    required 
                    label="Name" 
                    placeholder="Enter name" 
                    value={createName} 
                    onChange={(e) => setCreateName(e.target.value)}
                />
                <TextField 
                    label="Description"
                    placeholder="Enter description"
                    multiline
                    value={createDescription}
                    onChange={(e) => setCreateDescription(e.target.value)}
                />
                <br />
                <TextField 
                    required
                    label="Select frequency"
                    value={createFrequency}
                    onChange={(e) => setCreateFrequency(e.target.value)}
                    select
                    fullWidth
                    sx={{ marginTop: 2 }}
                > 
                    <MenuItem value="daily">daily</MenuItem>
                    <MenuItem value="weekly">weekly</MenuItem>
                    <MenuItem value="monthly">monthly</MenuItem>
                </TextField>
                <TextField
                    required
                    label="Enter Remind Time"
                    type="time"
                    value={createRemindTime}
                    onChange={(e) => setCreateRemindTime(e.target.value)}
                    sx={{ marginTop: 2, width: "100%" }}
                />
                <TextField 
                    label="Enter timezone"
                    placeholder="Timezone" 
                    value={createTimezone} 
                    onChange={(e) => setCreateTimezone(e.target.value)}
                    sx={{ marginTop: 2, marginBottom: 2, width: "100%" }}
                />
                <Button loading={isCreating} type="submit" variant="outlined">Create habit!</Button>
            </Box>
            {updateHabitId && (
                <Box component="form" onSubmit={updateHabit} sx={{ marginBottom: 2 }}>
                    <Typography variant="h3">Udpate habit { updateHabitId }</Typography>
                    <TextField 
                        required 
                        label="Name" 
                        placeholder="Enter name" 
                        value={updateName} 
                        onChange={(e) => setUpdateName(e.target.value)}
                    />
                    <TextField 
                        label="Description"
                        placeholder="Enter description"
                        multiline
                        value={updateDescription}
                        onChange={(e) => setUpdateDescription(e.target.value)}
                    />
                    <br />
                    <TextField 
                        required
                        label="Select frequency"
                        value={updateFrequency}
                        onChange={(e) => setUpdateFrequency(e.target.value)}
                        select
                        fullWidth
                        sx={{ marginTop: 2 }}
                    > 
                        <MenuItem value="daily">daily</MenuItem>
                        <MenuItem value="weekly">weekly</MenuItem>
                        <MenuItem value="monthly">monthly</MenuItem>
                    </TextField>
                    <TextField
                        required
                        label="Enter Remind Time"
                        type="time"
                        value={updateRemindTime.split('.').slice(0,1).join('.')}
                        onChange={(e) => setUpdateRemindTime(e.target.value)}
                        sx={{ marginTop: 2, width: "100%" }}
                    />
                    <TextField 
                        label="Enter timezone"
                        placeholder="Timezone" 
                        value={updateTimezone} 
                        onChange={(e) => setUpdateTimezone(e.target.value)}
                        sx={{ marginTop: 2, marginBottom: 2, width: "100%" }}
                    />
                    <Button loading={isUpdating} type="submit" variant="outlined">Update habit!</Button>
                </Box>
            )}
            <Box>
                {habits?.map((habit: Habit) => (
                    <Card variant="outlined" key={habit.id}>
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
                        <IconButton
                            loading={isUpdating}
                            aria-label="update"
                            onClick={() => {
                                setUpdateHabitId(String(habit.id))
                                setUpdateUserId(String(habit.user_id))
                                setUpdateName(habit.name)
                                setUpdateDescription(habit.description || "")
                                setUpdateRemindTime(habit.remind_time)
                                setUpdateFrequency(habit.frequency)
                                setUpdateTimezone(habit.timezone)
                            }}
                        >
                            <Update />
                        </IconButton>
                        <IconButton 
                            loading={isDeleting}
                            aria-label="delete"
                            onClick={() => {
                                setDeletingId(String(habit.id))
                                deleteHabit()
                            }}
                        >
                            <Delete />
                        </IconButton>
                    </Card>
                ))}
            </Box>
        </Box>
    )
}

export default HomePage;