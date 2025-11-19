# AionMCP Agent Examples

This document provides practical examples of how to create agents that use AionMCP tools.

## Example 1: Simple Data Aggregator Agent

This agent fetches data from multiple APIs and combines the results.

```python
#!/usr/bin/env python3
"""
Simple agent that uses AionMCP tools to aggregate data from multiple sources
"""
import json
import requests
from typing import List, Dict, Any

class AionMCPAgent:
    def __init__(self, server_url: str = "http://localhost:8080"):
        self.server_url = server_url
        self.api_url = f"{server_url}/api/v1"
        
    def get_tools(self) -> List[Dict[str, Any]]:
        """Fetch available tools from AionMCP server"""
        response = requests.get(f"{self.api_url}/tools")
        if response.status_code == 200:
            return response.json().get("tools", [])
        return []
    
    def execute_tool(self, tool_name: str, params: Dict[str, Any]) -> Any:
        """Execute a tool on the AionMCP server"""
        response = requests.post(
            f"{self.api_url}/tools/{tool_name}/execute",
            json=params,
            headers={"Content-Type": "application/json"}
        )
        if response.status_code == 200:
            return response.json()
        raise Exception(f"Tool execution failed: {response.text}")
    
    def get_user_data(self, user_id: int) -> Dict[str, Any]:
        """Fetch user and their posts using AionMCP tools"""
        # Get user details
        user = self.execute_tool("getUser", {"userId": user_id})
        
        # Get user's posts
        posts = self.execute_tool("listPosts", {})
        user_posts = [p for p in posts if p.get("userId") == user_id]
        
        return {
            "user": user,
            "posts": user_posts,
            "post_count": len(user_posts)
        }
    
    def get_country_info(self, country_name: str) -> Dict[str, Any]:
        """Fetch country information"""
        countries = self.execute_tool(
            "getCountriesByName",
            {"name": country_name}
        )
        return countries[0] if countries else None
    
    def find_pokemon(self, pokemon_name: str) -> Dict[str, Any]:
        """Find Pokémon by name"""
        return self.execute_tool(
            "getPokemon",
            {"id": pokemon_name.lower()}
        )


# Example usage
if __name__ == "__main__":
    agent = AionMCPAgent()
    
    # List available tools
    print("Available tools:")
    tools = agent.get_tools()
    for tool in tools[:5]:
        print(f"  - {tool['name']}: {tool['description']}")
    
    # Get user and posts
    print("\nFetching user 1 and posts...")
    user_data = agent.get_user_data(1)
    print(f"User: {user_data['user']['name']}")
    print(f"Posts: {user_data['post_count']}")
    
    # Get country info
    print("\nFetching country information...")
    country = agent.get_country_info("Japan")
    if country:
        print(f"Country: {country['name']['common']}")
        print(f"Capital: {country['capital']}")
    
    # Get Pokémon
    print("\nFetching Pokémon...")
    pokemon = agent.find_pokemon("pikachu")
    print(f"Name: {pokemon['name'].title()}")
    print(f"Types: {[t['type']['name'] for t in pokemon['types']]}")
```

## Example 2: Research Agent

An agent that researches topics by aggregating data from multiple sources.

```python
#!/usr/bin/env python3
"""
Research agent that gathers information from multiple APIs
"""
import json
from aionmcp_agent import AionMCPAgent

class ResearchAgent(AionMCPAgent):
    """Agent focused on research and information gathering"""
    
    def research_country(self, country_name: str) -> Dict[str, Any]:
        """Research a country comprehensively"""
        print(f"Researching {country_name}...")
        
        # Get country info
        country = self.execute_tool(
            "getCountriesByName",
            {"name": country_name}
        )
        
        if not country:
            return {"error": f"Country {country_name} not found"}
        
        country_data = country[0]
        
        research = {
            "country": country_data['name']['common'],
            "official_name": country_data['name']['official'],
            "region": country_data.get('region', 'Unknown'),
            "subregion": country_data.get('subregion', 'Unknown'),
            "capital": country_data.get('capital', ['Unknown'])[0] if country_data.get('capital') else 'Unknown',
            "population": country_data.get('population', 'Unknown'),
            "area": country_data.get('area', 'Unknown'),
            "timezones": country_data.get('timezones', []),
            "currencies": country_data.get('currencies', {}),
            "languages": country_data.get('languages', {}),
        }
        
        return research
    
    def find_related_countries(self, region: str) -> List[str]:
        """Find all countries in a region"""
        countries = self.execute_tool(
            "getCountriesByRegion",
            {"region": region}
        )
        
        return [c['name']['common'] for c in countries]
    
    def explore_pokemon_types(self, pokemon_name: str) -> Dict[str, Any]:
        """Explore a Pokémon and its type characteristics"""
        pokemon = self.execute_tool(
            "getPokemon",
            {"id": pokemon_name.lower()}
        )
        
        return {
            "name": pokemon['name'].title(),
            "height": pokemon['height'] / 10,  # Convert to meters
            "weight": pokemon['weight'] / 10,  # Convert to kg
            "types": [t['type']['name'] for t in pokemon['types']],
            "abilities": [a['ability']['name'] for a in pokemon['abilities']],
            "base_experience": pokemon['base_experience'],
        }


# Example usage
if __name__ == "__main__":
    agent = ResearchAgent()
    
    # Research a country
    print("=== Country Research ===")
    research = agent.research_country("Japan")
    print(json.dumps(research, indent=2))
    
    # Find related countries
    print("\n=== Countries in Asia ===")
    asian_countries = agent.find_related_countries("Asia")
    for country in asian_countries[:5]:
        print(f"  - {country}")
    
    # Explore Pokémon
    print("\n=== Pokémon Analysis ===")
    poke_info = agent.explore_pokemon_types("charizard")
    print(json.dumps(poke_info, indent=2))
```

## Example 3: Learning from User Interactions

An agent that learns from tool execution patterns.

```python
#!/usr/bin/env python3
"""
Agent that learns from user interactions and adapts behavior
"""
import time
import json
from datetime import datetime
from aionmcp_agent import AionMCPAgent

class LearningAgent(AionMCPAgent):
    """Agent that learns from interactions"""
    
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.interaction_history = []
        self.learned_patterns = {}
    
    def track_tool_execution(self, tool_name: str, params: Dict, result: Any, execution_time: float):
        """Track tool execution for learning"""
        interaction = {
            "timestamp": datetime.now().isoformat(),
            "tool": tool_name,
            "parameters": params,
            "success": result is not None,
            "execution_time": execution_time,
            "result_size": len(str(result))
        }
        self.interaction_history.append(interaction)
        
        # Learn from pattern
        if tool_name not in self.learned_patterns:
            self.learned_patterns[tool_name] = {
                "total_calls": 0,
                "avg_execution_time": 0,
                "success_rate": 0,
                "common_params": {}
            }
        
        pattern = self.learned_patterns[tool_name]
        pattern["total_calls"] += 1
        pattern["avg_execution_time"] = (
            (pattern["avg_execution_time"] * (pattern["total_calls"] - 1) + execution_time) /
            pattern["total_calls"]
        )
    
    def smart_execute(self, tool_name: str, params: Dict) -> Any:
        """Execute tool and learn from it"""
        start_time = time.time()
        try:
            result = self.execute_tool(tool_name, params)
            execution_time = time.time() - start_time
            self.track_tool_execution(tool_name, params, result, execution_time)
            return result
        except Exception as e:
            execution_time = time.time() - start_time
            self.track_tool_execution(tool_name, params, None, execution_time)
            raise
    
    def get_insights(self) -> Dict[str, Any]:
        """Analyze learned patterns"""
        insights = {
            "total_interactions": len(self.interaction_history),
            "tools_used": list(self.learned_patterns.keys()),
            "slowest_tools": sorted(
                [(name, pattern["avg_execution_time"]) 
                 for name, pattern in self.learned_patterns.items()],
                key=lambda x: x[1],
                reverse=True
            )[:5],
            "most_used_tools": sorted(
                [(name, pattern["total_calls"]) 
                 for name, pattern in self.learned_patterns.items()],
                key=lambda x: x[1],
                reverse=True
            )[:5],
        }
        return insights
    
    def save_learning_data(self, filepath: str):
        """Save learned data to file"""
        learning_data = {
            "interaction_history": self.interaction_history,
            "learned_patterns": self.learned_patterns,
            "insights": self.get_insights(),
        }
        with open(filepath, 'w') as f:
            json.dump(learning_data, f, indent=2)
        print(f"Learning data saved to {filepath}")


# Example usage
if __name__ == "__main__":
    agent = LearningAgent()
    
    # Perform various operations and learn
    print("Learning from tool executions...")
    
    # Execute various tools
    tools_to_test = [
        ("listPosts", {}),
        ("listUsers", {}),
        ("getUser", {"userId": 1}),
        ("listComments", {"postId": 1}),
    ]
    
    for tool_name, params in tools_to_test:
        try:
            result = agent.smart_execute(tool_name, params)
            print(f"✓ {tool_name}: Success")
        except Exception as e:
            print(f"✗ {tool_name}: {e}")
    
    # Display insights
    print("\n=== Agent Insights ===")
    insights = agent.get_insights()
    print(json.dumps(insights, indent=2))
    
    # Save learning data
    agent.save_learning_data("agent_learning.json")
```

## Example 4: Multi-API Orchestration Agent

An agent that coordinates between multiple APIs to accomplish complex tasks.

```python
#!/usr/bin/env python3
"""
Agent that orchestrates multiple APIs to complete complex workflows
"""
from aionmcp_agent import AionMCPAgent

class OrchestrationAgent(AionMCPAgent):
    """Agent that coordinates multiple tool calls"""
    
    def create_user_profile_report(self, user_id: int) -> str:
        """Create a comprehensive report for a user"""
        # Fetch user
        user = self.execute_tool("getUser", {"userId": user_id})
        
        # Fetch user's posts
        posts = self.execute_tool("listPosts", {})
        user_posts = [p for p in posts if p.get("userId") == user_id]
        
        # Create report
        report = f"""
=== User Profile Report ===
Name: {user['name']}
Email: {user['email']}
Phone: {user['phone']}
Website: {user['website']}
Company: {user['company']['name']}

Total Posts: {len(user_posts)}
Recent Posts:
"""
        for post in user_posts[-3:]:
            report += f"\n- {post['title']}"
        
        return report
    
    def generate_global_statistics(self) -> Dict[str, Any]:
        """Generate statistics from all available data"""
        # Get all posts
        posts = self.execute_tool("listPosts", {})
        
        # Get all users
        users = self.execute_tool("listUsers", {})
        
        # Calculate stats
        stats = {
            "total_users": len(users),
            "total_posts": len(posts),
            "avg_posts_per_user": len(posts) / len(users) if users else 0,
            "users_by_city": self._group_by_city(users),
        }
        
        return stats
    
    def _group_by_city(self, users: List[Dict]) -> Dict[str, int]:
        """Group users by city"""
        cities = {}
        for user in users:
            city = user.get('address', {}).get('city', 'Unknown')
            cities[city] = cities.get(city, 0) + 1
        return cities


# Example usage
if __name__ == "__main__":
    agent = OrchestrationAgent()
    
    # Generate user report
    print(agent.create_user_profile_report(1))
    
    # Generate global statistics
    print("\n=== Global Statistics ===")
    stats = agent.generate_global_statistics()
    for key, value in stats.items():
        print(f"{key}: {value}")
```

## Running These Examples

### Prerequisites

```bash
# Install required packages
pip install requests

# Start AionMCP server
cd /Users/kiran/Documents/GitHub/aionmcp
./server
```

### Run an example

```bash
# Save one of the above scripts as agent.py
python3 agent.py
```

## Integration with AionMCP Learning Engine

These agents will automatically generate learning data that AionMCP can analyze:

1. **Execution patterns** are recorded in the learning database
2. **Tool usage statistics** are accumulated
3. **Reflection documents** are generated in `docs/reflections/`
4. **Patterns** help optimize future tool selection

## Advanced Features to Implement

1. **Caching**: Cache results from slow tools
2. **Retry Logic**: Automatically retry failed operations
3. **Error Recovery**: Handle API failures gracefully
4. **Rate Limiting**: Respect API limits
5. **Parallel Execution**: Execute independent tools concurrently
6. **Context Preservation**: Maintain state across tool calls
7. **Feedback Loop**: Use learning data to improve decisions

## References

- AionMCP API: `docs/api.md`
- Tool Registry: See `internal/core/registry.go`
- Learning Engine: See `internal/selflearn/engine.go`
